package broker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/knita-io/knita/api/broker/v1"
	brokerv1 "github.com/knita-io/knita/api/executor/v1"
)

// LocalBroker brokers runtimes via the local mediator.
// The Broker's job is to locate executors (via the mediator API) that are capable of bidding on the Tender
type LocalBroker struct {
	v1.UnimplementedRuntimeBrokerServer
	log        *zap.SugaredLogger
	socketPath string
}

func NewLocalBroker(log *zap.SugaredLogger, socketPath string) *LocalBroker {
	return &LocalBroker{log: log.Named("local_broker"), socketPath: socketPath}
}

func (b *LocalBroker) Tender(ctx context.Context, req *v1.RuntimeTender) (*v1.RuntimeContracts, error) {
	log := b.log.With("tender_id", req.TenderId)
	log.Infow("Brokering runtime contract...")
	client, done, err := b.dialLocalExecutor()
	if err != nil {
		return nil, err
	}
	defer done()
	intro, err := client.Introspect(ctx, &brokerv1.IntrospectRequest{})
	if err != nil {
		return nil, fmt.Errorf("error introspecting local mediator: %w", err)
	}
	var contracts []*v1.RuntimeContract
	if b.canBid(intro, req) {
		contracts = append(contracts, &v1.RuntimeContract{
			ContractId: uuid.New().String(),
			RuntimeId:  uuid.New().String(),
			Opts:       req.Opts,
		})
	}
	log.Infow("Brokered contracts", "n_contracts", len(contracts))
	return &v1.RuntimeContracts{Contracts: contracts}, nil
}

func (b *LocalBroker) Settle(ctx context.Context, req *v1.RuntimeContract) (*v1.RuntimeSettlement, error) {
	log := b.log.With("contract_id", req.ContractId)
	log.Infow("Settling contract...")
	// TODO: Hook up auth
	log.Infow("Settled contract")
	return &v1.RuntimeSettlement{
		ConnectionInfo: &v1.RuntimeConnectionInfo{
			Transport: &v1.RuntimeConnectionInfo_Unix{
				Unix: &v1.RuntimeTransportUnix{SocketPath: b.socketPath},
			}},
	}, nil
}

func (b *LocalBroker) canBid(intro *brokerv1.IntrospectResponse, tender *v1.RuntimeTender) bool {
	return isSubset(tender.Opts.Labels, intro.Labels)
}

func (b *LocalBroker) dialLocalExecutor() (brokerv1.ExecutorClient, func(), error) {
	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}
	conn, err := grpc.Dial(b.socketPath, grpc.WithInsecure(), grpc.WithDialer(dialer), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing mediator: %w", err)
	}
	return brokerv1.NewExecutorClient(conn), func() { conn.Close() }, nil
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
