package server

import (
	"context"
	"reflect"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func MakeUnaryServerLogInterceptor(syslog *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			syslog.Warnw("Recv message failed", "req", reflect.TypeOf(req).String(), "res", reflect.TypeOf(resp).String(), "method", info.FullMethod, "err", err.Error())
		} else {
			syslog.Debugw("Received message", "req", reflect.TypeOf(req).String(), "res", reflect.TypeOf(resp).String(), "method", info.FullMethod)
		}
		return resp, err
	}
}
