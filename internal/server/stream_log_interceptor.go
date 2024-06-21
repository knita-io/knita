package server

import (
	"errors"
	"io"
	"reflect"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type wrapper struct {
	grpc.ServerStream
	syslog *zap.SugaredLogger
	id     string
}

func (w *wrapper) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			w.syslog.Warnw("Recv message failed", "stream_id", w.id, "msg", reflect.TypeOf(m).String(), "err", err.Error())
		}
	} else {
		w.syslog.Debugw("Received message", "stream_id", w.id, "msg", reflect.TypeOf(m).String())
	}
	return err
}

func (w *wrapper) SendMsg(m any) error {
	err := w.ServerStream.SendMsg(m)
	if err != nil {
		w.syslog.Warnw("Send message failed", "stream_id", w.id, "msg", reflect.TypeOf(m).String(), "err", err.Error())
	} else {
		w.syslog.Debugw("Sent message", "stream_id", w.id, "msg", reflect.TypeOf(m).String())
	}
	return err
}

func MakeStreamServerLogInterceptor(syslog *zap.SugaredLogger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id := xid.New().String()
		syslog.Infow("Stream opened", "stream_id", id, "method", info.FullMethod)
		err := handler(srv, &wrapper{ServerStream: ss, syslog: syslog, id: id})
		if err != nil {
			syslog.Warnw("Stream closed", "stream_id", id, "err", err.Error())
		} else {
			syslog.Infow("Stream closed", "stream_id", id)
		}
		return err
	}
}
