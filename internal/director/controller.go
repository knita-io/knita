package director

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	brokerv1 "github.com/knita-io/knita/api/broker/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/file"
)

type BuildController struct {
	syslog      *zap.SugaredLogger
	log         *Log
	buildID     string
	broker      brokerv1.RuntimeBrokerClient
	localWorkFS file.WriteFS
}

func NewBuildController(syslog *zap.SugaredLogger, log *Log, buildID string, broker brokerv1.RuntimeBrokerClient, localWorkFS file.WriteFS) *BuildController {
	return &BuildController{
		syslog:      syslog.Named("director"),
		log:         log,
		buildID:     buildID,
		broker:      broker,
		localWorkFS: localWorkFS,
	}
}

func (c *BuildController) BuildID() string {
	return c.buildID
}

func (c *BuildController) Log() *Log {
	return c.log
}

func (c *BuildController) Runtime(ctx context.Context, opts *executorv1.Opts) (*Runtime, error) {
	log := c.syslog
	log.Infow("Requesting runtime from broker...", "opts", opts)
	tenderID := uuid.New().String()
	tender := &brokerv1.RuntimeTender{TenderId: tenderID, Opts: opts}
	c.log.Publish(executorv1.NewRuntimeTenderStartEvent(tenderID, opts))
	runtimeRes, err := c.broker.Tender(ctx, tender)
	if err != nil {
		return nil, fmt.Errorf("error brokering runtime: %w", err)
	}
	c.log.Publish(executorv1.NewRuntimeTenderEndEvent(tenderID))
	if len(runtimeRes.Contracts) == 0 {
		return nil, fmt.Errorf("error no runtime contracts received; unable to locate suitable executor to host runtime")
	}
	contract := runtimeRes.Contracts[0]
	log.Infow("Selected runtime contract", "contract_id", contract.ContractId)
	rid := contract.RuntimeId
	c.log.Publish(executorv1.NewRuntimeSettlementStartEvent(tenderID, contract.ContractId, rid))
	settlementRes, err := c.broker.Settle(ctx, contract)
	if err != nil {
		return nil, fmt.Errorf("error settling contract: %w", err)
	}
	c.log.Publish(executorv1.NewRuntimeSettlementEndEvent(tenderID, contract.ContractId, rid))
	log.Infow("Settled runtime contract", "contract_id", contract.ContractId)
	report := c.makeSelectionReport(tender, runtimeRes.Contracts, contract, settlementRes)
	c.log.Printf(report)
	log = log.With("runtime_id", rid)
	var rClient executorv1.ExecutorClient
	switch trans := settlementRes.ConnectionInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}
		conn, err := grpc.DialContext(ctx, trans.Unix.SocketPath, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("error dialing executor via unix domain socket: %w", err)
		}
		rClient = executorv1.NewExecutorClient(conn)
	case *brokerv1.RuntimeConnectionInfo_Tcp:
		conn, err := grpc.DialContext(ctx, trans.Tcp.Address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("error dialing executor via tcp: %w", err)
		}
		rClient = executorv1.NewExecutorClient(conn)
	default:
		return nil, fmt.Errorf("error unknown transport: %T", trans)
	}
	log.Info("Connected to executor")
	r := newRuntime(c.syslog, c.log, c.buildID, rid, rClient, c.localWorkFS, contract.Opts)
	err = r.start(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating runtime: %w", err)
	}
	return r, nil
}

// makeSelectionReport generates a concise text report about the runtime tender process
// and results, suitable for inclusion in the build log.
func (c *BuildController) makeSelectionReport(
	tender *brokerv1.RuntimeTender,
	contracts []*brokerv1.RuntimeContract,
	selectedContractOrNil *brokerv1.RuntimeContract,
	settlement *brokerv1.RuntimeSettlement) string {

	name := tender.TenderId
	if tName, ok := tender.Opts.Tags["name"]; ok {
		name = tName
	}
	requires := strings.Join(tender.Opts.Labels, ",")
	output := fmt.Sprintf("Elegible Executors for Runtime: %s (type=%s, requires=%s)\n", name, tender.Opts.Type, requires)

	for _, contract := range contracts {
		output += fmt.Sprintf("  %s (os=%s, arch=%s, cpu=%d, memory=%d)\n", contract.ExecutorInfo.Name, contract.SysInfo.Os, contract.SysInfo.Arch, contract.SysInfo.TotalCpuCores, contract.SysInfo.TotalMemory)
	}

	if selectedContractOrNil == nil {
		output += "Error no executors available"
	} else {
		connInfo := ""
		switch transport := settlement.ConnectionInfo.Transport.(type) {
		case *brokerv1.RuntimeConnectionInfo_Unix:
			connInfo = fmt.Sprintf(" (unix:/%s)", transport.Unix.SocketPath)
		case *brokerv1.RuntimeConnectionInfo_Tcp:
			connInfo = fmt.Sprintf(" (tcp://%s)", transport.Tcp.Address)
		}
		output += fmt.Sprintf("Selected Executor: %s%s", selectedContractOrNil.ExecutorInfo.Name, connInfo)
	}
	return output
}
