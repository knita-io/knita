package director

import (
	"context"
	"fmt"
	"net"
	stdruntime "runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/pbnjay/memory"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	brokerv1 "github.com/knita-io/knita/api/broker/v1"
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/file"
	"github.com/knita-io/knita/internal/version"
)

type Build struct {
	syslog      *zap.SugaredLogger
	log         *Log
	buildID     string
	broker      brokerv1.RuntimeBrokerClient
	localWorkFS file.WriteFS
}

func NewBuild(syslog *zap.SugaredLogger, log *Log, buildID string, broker brokerv1.RuntimeBrokerClient, localWorkFS file.WriteFS) *Build {
	return &Build{
		syslog:      syslog.Named("director"),
		log:         log,
		buildID:     buildID,
		broker:      broker,
		localWorkFS: localWorkFS,
	}
}

// BuildID returns the build ID associated with the Build.
func (c *Build) BuildID() string {
	return c.buildID
}

// Log returns the log associated with the Build.
func (c *Build) Log() *Log {
	return c.log
}

// Run executes fn, wrapping it in build start and end events.
func (c *Build) Run(fn func() error) error {
	info := &builtinv1.DirectorInfo{
		Version: version.Version,
		SysInfo: &executorv1.SystemInfo{
			Os:            stdruntime.GOOS,
			Arch:          stdruntime.GOARCH,
			TotalCpuCores: uint32(stdruntime.NumCPU()),
			TotalMemory:   memory.TotalMemory(),
		},
	}
	c.log.Publish(&builtinv1.BuildStartEvent{BuildId: c.buildID, DirectorInfo: info})
	return WithEndEvent(fn, func() {
		c.log.Publish(&builtinv1.BuildEndEvent{BuildId: c.buildID,
			Status: &builtinv1.BuildEndEvent_Result{Result: &builtinv1.BuildResult{}}})
	}, func(err error) {
		c.log.Publish(&builtinv1.BuildEndEvent{BuildId: c.buildID,
			Status: &builtinv1.BuildEndEvent_Error{Error: &builtinv1.Error{Message: err.Error()}}})
	})
}

// OpenRuntime requests a runtime from the broker configured with the given options.
func (c *Build) OpenRuntime(ctx context.Context, opts *executorv1.Opts) (*Runtime, error) {
	c.syslog.Infow("Tendering runtime...", "opts", opts)
	runtimeRes, err := c.tenderRuntime(ctx, opts)
	if err != nil {
		return nil, err
	}
	contract := runtimeRes.Contracts[0]
	c.syslog.Infow("Selected runtime contract", "contract_id", contract.ContractId)
	settlementRes, err := c.settleRuntime(ctx, contract)
	if err != nil {
		return nil, err
	}
	c.syslog.Infow("Settled runtime contract", "contract_id", contract.ContractId)
	c.log.Printf(c.makeSelectionReport(runtimeRes.Contracts, contract, settlementRes))
	rClient, err := c.makeExecutorClient(ctx, settlementRes.ConnectionInfo)
	if err != nil {
		return nil, err
	}
	c.syslog.Info("Connected to executor")
	r := newRuntime(c.syslog, c.log, c.buildID, contract.RuntimeId, rClient, c.localWorkFS, contract.Opts)
	err = r.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating runtime: %w", err)
	}
	return r, nil
}

// tenderRuntime puts a runtime out for tender.
// Returns an error if no contracts were received.
func (c *Build) tenderRuntime(ctx context.Context, opts *executorv1.Opts) (*brokerv1.TenderResponse, error) {
	tenderID := uuid.New().String()
	c.log.Publish(&builtinv1.RuntimeTenderStartEvent{BuildId: c.buildID, TenderId: tenderID, Opts: opts})
	return WithUnaryEndEvent(func() (*brokerv1.TenderResponse, error) {
		tender := &brokerv1.TenderRequest{BuildId: c.buildID, TenderId: tenderID, Opts: opts}
		res, err := c.broker.Tender(ctx, tender)
		if err != nil {
			return nil, err
		}
		if len(res.Contracts) == 0 {
			return nil, fmt.Errorf("error no runtime contracts received; unable to locate an executor to host the runtime")
		}
		return res, nil
	}, func(res *brokerv1.TenderResponse) {
		c.log.Publish(&builtinv1.RuntimeTenderEndEvent{TenderId: tenderID,
			Status: &builtinv1.RuntimeTenderEndEvent_Result{Result: &builtinv1.RuntimeTenderResult{Contracts: res.Contracts}}})
	}, func(err error) {
		c.log.Publish(&builtinv1.RuntimeTenderEndEvent{TenderId: tenderID,
			Status: &builtinv1.RuntimeTenderEndEvent_Error{Error: &builtinv1.Error{Message: err.Error()}}})
	})
}

// settleRuntime settles a runtime contract that was previously tendered and returns the results.
func (c *Build) settleRuntime(ctx context.Context, contract *brokerv1.RuntimeContract) (*brokerv1.SettlementResponse, error) {
	c.log.Publish(&builtinv1.RuntimeSettlementStartEvent{TenderId: contract.TenderId, ContractId: contract.ContractId, RuntimeId: contract.RuntimeId})
	return WithUnaryEndEvent(func() (*brokerv1.SettlementResponse, error) {
		settlementRes, err := c.broker.Settle(ctx, &brokerv1.SettlementRequest{Contract: contract})
		if err != nil {
			return nil, fmt.Errorf("error settling contract: %w", err)
		}
		return settlementRes, nil
	}, func(res *brokerv1.SettlementResponse) {
		c.log.Publish(&builtinv1.RuntimeSettlementEndEvent{TenderId: contract.TenderId, ContractId: contract.ContractId, RuntimeId: contract.RuntimeId,
			Status: &builtinv1.RuntimeSettlementEndEvent_Result{Result: &builtinv1.RuntimeSettlementResult{}}})
	}, func(err error) {
		c.log.Publish(&builtinv1.RuntimeSettlementEndEvent{TenderId: contract.TenderId, ContractId: contract.ContractId, RuntimeId: contract.RuntimeId,
			Status: &builtinv1.RuntimeSettlementEndEvent_Error{Error: &builtinv1.Error{Message: err.Error()}}})
	})
}

// makeExecutorClient returns an executor client configured to connect to the executor in connInfo.
func (c *Build) makeExecutorClient(ctx context.Context, connInfo *brokerv1.RuntimeConnectionInfo) (executorv1.ExecutorClient, error) {
	switch trans := connInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}
		conn, err := grpc.DialContext(ctx, trans.Unix.SocketPath, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("error dialing executor via unix domain socket: %w", err)
		}
		return executorv1.NewExecutorClient(conn), nil
	case *brokerv1.RuntimeConnectionInfo_Tcp:
		conn, err := grpc.DialContext(ctx, trans.Tcp.Address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("error dialing executor via tcp: %w", err)
		}
		return executorv1.NewExecutorClient(conn), nil
	default:
		return nil, fmt.Errorf("error unknown transport: %T", trans)
	}
}

// makeSelectionReport generates a concise text report about the runtime tender process
// and results, suitable for inclusion in the build log.
func (c *Build) makeSelectionReport(
	contracts []*brokerv1.RuntimeContract,
	selectedContract *brokerv1.RuntimeContract, settlement *brokerv1.SettlementResponse) string {

	name := selectedContract.TenderId
	if tName, ok := selectedContract.Opts.Tags["name"]; ok {
		name = tName
	}
	requires := strings.Join(selectedContract.Opts.Labels, ",")
	output := fmt.Sprintf("Elegible Executors for Runtime: %s (type=%s, requires=%s)\n", name, selectedContract.Opts.Type, requires)
	for _, contract := range contracts {
		output += fmt.Sprintf("  %s (os=%s, arch=%s, cpu=%d, memory=%d)\n", contract.ExecutorInfo.Name, contract.SysInfo.Os, contract.SysInfo.Arch, contract.SysInfo.TotalCpuCores, contract.SysInfo.TotalMemory)
	}
	connInfo := ""
	switch transport := settlement.ConnectionInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		connInfo = fmt.Sprintf(" (unix:/%s)", transport.Unix.SocketPath)
	case *brokerv1.RuntimeConnectionInfo_Tcp:
		connInfo = fmt.Sprintf(" (tcp://%s)", transport.Tcp.Address)
	}
	output += fmt.Sprintf("Selected Executor: %s%s", selectedContract.ExecutorInfo.Name, connInfo)
	return output
}
