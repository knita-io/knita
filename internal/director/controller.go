package director

import (
	"context"
	"fmt"
	"net"
	"time"

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

func (c *BuildController) Runtime(ctx context.Context, opts *executorv1.Opts) (*Runtime, error) {
	log := c.syslog
	log.Infow("Requesting runtime from broker...", "opts", opts)
	runtimeRes, err := c.broker.Tender(ctx, &brokerv1.RuntimeTender{
		TenderId: uuid.New().String(),
		Opts:     opts,
	})
	if err != nil {
		return nil, fmt.Errorf("error brokering runtime: %w", err)
	}
	if len(runtimeRes.Contracts) == 0 {
		return nil, fmt.Errorf("error no runtime contracts received; unable to locate suitable executor to host runtime")
	}
	contract := runtimeRes.Contracts[0]
	rid := contract.RuntimeId
	log.Infow("Selected runtime contract", "contract_id", contract.ContractId)

	settlementRes, err := c.broker.Settle(ctx, contract)
	if err != nil {
		return nil, fmt.Errorf("error settling contract: %w", err)
	}
	log.Infow("Settled runtime contract", "contract_id", contract.ContractId)
	log = log.With("runtime_id", rid)

	var rClient executorv1.ExecutorClient
	switch trans := settlementRes.ConnectionInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		dialer := func(addr string, t time.Duration) (net.Conn, error) {
			return net.Dial("unix", addr)
		}
		conn, err := grpc.Dial(trans.Unix.SocketPath, grpc.WithInsecure(), grpc.WithDialer(dialer), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("error dialing mediator: %w", err)
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
