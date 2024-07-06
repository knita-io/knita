package fixed

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
	Connection *brokerv1.RuntimeConnectionInfo
}

type executorState struct {
	id            string
	config        *ExecutorConfig
	client        executorv1.ExecutorClient
	done          func()
	introspection *executorv1.IntrospectResponse
}

// Server brokers runtimes across a fixed (at run time) set of well-known remote executors.
type Server struct {
	brokerv1.UnimplementedRuntimeBrokerServer
	syslog        *zap.SugaredLogger
	config        Config
	initOnce      sync.Once
	executorsByID map[string]*executorState
}

// NewServer creates a new instance of the Server struct with the provided logger and config.
func NewServer(syslog *zap.SugaredLogger, config Config) *Server {
	return &Server{
		syslog:        syslog.Named("fixed_broker"),
		config:        config,
		executorsByID: make(map[string]*executorState),
	}
}

// Tender brokers a runtime contract based on the provided runtime tender.
func (b *Server) Tender(ctx context.Context, req *brokerv1.RuntimeTender) (*brokerv1.RuntimeContracts, error) {
	if err := validateRuntimeTender(req); err != nil {
		return nil, err
	}
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
				SysInfo:      executor.introspection.SysInfo,
				ExecutorInfo: executor.introspection.ExecutorInfo,
			})
		}
	}
	b.syslog.Infow("Brokered contracts", "n_contracts", len(contracts))
	return &brokerv1.RuntimeContracts{Contracts: contracts}, nil
}

// Settle settles the contract identified by the provided runtime contract.
func (b *Server) Settle(ctx context.Context, req *brokerv1.RuntimeContract) (*brokerv1.RuntimeSettlement, error) {
	if err := validateRuntimeContract(req); err != nil {
		return nil, err
	}
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

func (b *Server) canBid(intro *executorv1.IntrospectResponse, tender *brokerv1.RuntimeTender) bool {
	return isSubset(tender.Opts.Labels, intro.Labels)
}

// init initializes the server by initializing each executor based on the provided configuration.
// If there is an error during initialization, it logs a warning message and continues.
// This method is intended to be called only once during the server's initialization process.
func (b *Server) init() {
	b.syslog.Infof("Initializing executors...")
	for _, execConfig := range b.config.Executors {
		state, err := b.initExecutor(execConfig)
		if err != nil {
			b.syslog.Warnf("Ignoring error initializing executor %q; Executor will be "+
				"unavailable to run builds: %v", b.connInfoToString(execConfig.Connection), err)
		} else {
			b.syslog.Infow("Initialized executor", "name", state.introspection.ExecutorInfo.Name)
		}
	}
}

// initExecutor initializes the executor based on the provided configuration.
// It dials the executor using the connection information, performs an introspection,
// and stores the executor state in the Server's executorsByID map.
// It returns an error if there is an issue dialing the executor or introspecting the executor.
func (b *Server) initExecutor(execConfig *ExecutorConfig) (*executorState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, done, err := b.dialExecutor(ctx, execConfig.Connection)
	if err != nil {
		return nil, fmt.Errorf("error dialing executor: %w", err)
	}
	introspect, err := client.Introspect(ctx, &executorv1.IntrospectRequest{})
	if err != nil {
		return nil, fmt.Errorf("error introspecting executor: %w", err)
	}
	id := uuid.New().String()
	state := &executorState{
		id:            id,
		config:        execConfig,
		client:        client,
		done:          done,
		introspection: introspect,
	}
	b.executorsByID[id] = state
	return state, nil
}

// dialExecutor dials the executor based on the provided connection information.
// It returns an ExecutorClient and a cleanup function to close the connection.
// It returns an error if there is an issue dialing the executor.
func (b *Server) dialExecutor(ctx context.Context, connInfo *brokerv1.RuntimeConnectionInfo) (executorv1.ExecutorClient, func(), error) {
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

// connInfoToString returns a human-readable string representation of the connection info.
func (b *Server) connInfoToString(connInfo *brokerv1.RuntimeConnectionInfo) string {
	switch t := connInfo.Transport.(type) {
	case *brokerv1.RuntimeConnectionInfo_Unix:
		return t.Unix.SocketPath
	case *brokerv1.RuntimeConnectionInfo_Tcp:
		return t.Tcp.Address
	default:
		return "unknown"
	}
}

// validateRuntimeTender validates the fields of a RuntimeTender request.
// It returns an error if any of the mandatory fields are empty, otherwise returns nil.
func validateRuntimeTender(req *brokerv1.RuntimeTender) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.TenderId == "" {
		return fmt.Errorf("empty tender_id")
	}
	if req.BuildId == "" {
		return fmt.Errorf("empty build_id")
	}
	if req.Opts == nil {
		return fmt.Errorf("empty opts")
	}
	// NOTE opts are not validated here as the broker client (and possibly the executor that
	// wins the tender) may be newer than the broker server. The broker server should only
	// validate inputs strictly needed to complete the tender and otherwise defer to the executor.
	return nil
}

// validateRuntimeContract validates the fields of a RuntimeContract request.
// It returns an error if any of the mandatory fields are empty, otherwise returns nil.
func validateRuntimeContract(req *brokerv1.RuntimeContract) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.TenderId == "" {
		return fmt.Errorf("empty tender_id")
	}
	if req.ContractId == "" {
		return fmt.Errorf("empty contract_id")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.Opts == nil {
		return fmt.Errorf("empty opts")
	}
	return nil
}

// subset returns true if the first array is completely
// contained in the second array. There must be at least
// the same number of duplicate values in second as there
// are in first.
func isSubset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value] += 1
	}
	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}
	return true
}
