package server

import (
	"context"
	"reflect"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func MakeUnaryServerLogInterceptor(syslog *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		id := xid.New().String()
		syslog.Debugw("Received message", "id", id, "req", reflect.TypeOf(req).String(), "method", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			syslog.Warnw("Sent response", "id", id, "res", reflect.TypeOf(resp).String(), "method", "err", err.Error(), info.FullMethod)
		} else {
			syslog.Debugw("Sent response", "id", id, "res", reflect.TypeOf(resp).String(), "method", info.FullMethod)
		}
		return resp, err
	}
}
