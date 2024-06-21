package broker

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	brokerv1 "github.com/knita-io/knita/api/broker/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type Config struct {
	Executors []*ExecutorConfig
}

type ExecutorConfig struct {
	Name       string
	Connection *brokerv1.RuntimeConnectionInfo
}

type executorState struct {
	id            string
	config        *ExecutorConfig
	client        executorv1.ExecutorClient
	done          func()
	introspection *executorv1.IntrospectResponse
}

// FixedBroker brokers runtimes across a fixed (at run time) set of well-known remote executors.
type FixedBroker struct {
	brokerv1.UnimplementedRuntimeBrokerServer
	syslog        *zap.SugaredLogger
	config        Config
	initOnce      sync.Once
	executorsByID map[string]*executorState
}

func NewFixedBroker(syslog *zap.SugaredLogger, config Config) *FixedBroker {
	return &FixedBroker{
		syslog:        syslog.Named("fixed_broker"),
		config:        config,
		executorsByID: make(map[string]*executorState),
	}
}

func (b *FixedBroker) Tender(ctx context.Context, req *brokerv1.RuntimeTender) (*brokerv1.RuntimeContracts, error) {
	b.initOnce.Do(b.init)
	syslog := b.syslog.With("tender_id", req.TenderId)
	syslog.Infow("Brokering runtime contract...")
	var contracts []*brokerv1.RuntimeContract
	for _, executor := range b.executorsByID {
		if b.canBid(executor.introspection, req) {
			// NOTE: Here we use the unique executor ID as the contract ID as all we really
			// care about is being able to map a future settlement request to an executor.
			contracts = append(contracts, &brokerv1.RuntimeContract{
				TenderId:   req.TenderId,
				ContractId: executor.id,
				RuntimeId:  uuid.New().String(),
				Opts:       req.Opts,
				// TODO: We don't currently enforce any resource limits on Docker containers, so it's
				//  accurate to just pass the executor host's sys info back as part of the contract, but
				//  eventually this will need to change.
				SysInfo: executor.introspection.SysInfo,
			})
		}
	}
	b.syslog.Infow("Brokered contracts", "n_contracts", len(contracts))
	return &brokerv1.RuntimeContracts{Contracts: contracts}, nil
}

func (b *FixedBroker) Settle(ctx context.Context, req *brokerv1.RuntimeContract) (*brokerv1.RuntimeSettlement, error) {
	syslog := b.syslog.With("contract_id", req.ContractId)
	syslog.Infow("Settling contract...")
	executor, ok := b.executorsByID[req.ContractId]
	if !ok {
		return nil, fmt.Errorf("executor not found")
	}
	// TODO: Hook up auth
	syslog.Infow("Settled contract")
	return &brokerv1.RuntimeSettlement{ConnectionInfo: executor.config.Connection}, nil
}

func (b *FixedBroker) canBid(intro *executorv1.IntrospectResponse, tender *brokerv1.RuntimeTender) bool {
	return isSubset(tender.Opts.Labels, intro.Labels)
}

func (b *FixedBroker) init() {
	b.syslog.Infof("Initializing executors...")
	for _, execConfig := range b.config.Executors {
		err := b.initExecutor(execConfig)
		if err != nil {
			b.syslog.Warnf("Ignoring error initializing executor %q; Executor will be "+
				"unavailable to run builds: %v", execConfig.Name, err)
		} else {
			b.syslog.Infow("Initialized executor", "name", execConfig.Name)
		}
	}
}

func (b *FixedBroker) initExecutor(execConfig *ExecutorConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, done, err := b.dialExecutor(ctx, execConfig.Connection)
	if err != nil {
		return fmt.Errorf("error dialing executor: %w", err)
	}
	introspect, err := client.Introspect(ctx, &executorv1.IntrospectRequest{})
	if err != nil {
		return fmt.Errorf("error introspecting executor: %w", err)
	}
	id := uuid.New().String()
	b.executorsByID[id] = &executorState{
		id:            id,
		config:        execConfig,
		client:        client,
		done:          done,
		introspection: introspect,
	}
	return nil
}

func (b *FixedBroker) dialExecutor(ctx context.Context, connInfo *brokerv1.RuntimeConnectionInfo) (executorv1.ExecutorClient, func(), error) {
	switch t := connInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}
		conn, err := grpc.DialContext(ctx, t.Unix.SocketPath, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithBlock())
		if err != nil {
			return nil, nil, fmt.Errorf("error dialing executor via unix domain socket: %w", err)
		}
		return executorv1.NewExecutorClient(conn), func() { conn.Close() }, nil
	case *brokerv1.RuntimeConnectionInfo_Tcp:
		conn, err := grpc.DialContext(ctx, t.Tcp.Address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, nil, fmt.Errorf("error dialing executor via tcp: %w", err)
		}
		return executorv1.NewExecutorClient(conn), func() { conn.Close() }, nil
	default:
		return nil, nil, fmt.Errorf("error unsupported connection type: %T", t)
	}
}
