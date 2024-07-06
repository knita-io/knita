package file

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type ReceiveState int

const (
	ReceiveStateAwaitingHeaders ReceiveState = iota
	ReceiveStateAwaitingBody
	ReceiveStateAwaitingTrailer
	ReceiveStateDone
)

type RecvCallback func(header *executorv1.FileTransferHeader)

type RecvOpt interface {
	Apply(*RecvOpts)
}

type RecvOpts struct {
	cb RecvCallback
}

type withRecvCallback struct {
	cb RecvCallback
}

func (o *withRecvCallback) Apply(opts *RecvOpts) {
	opts.cb = o.cb
}

func WithRecvCallback(cb RecvCallback) RecvOpt {
	return &withRecvCallback{cb: cb}
}

type Receiver struct {
	syslog *zap.SugaredLogger
	opts   *RecvOpts
	fs     WriteFS
	state  ReceiveState
	fh     File
	header *executorv1.FileTransferHeader
}

func NewReceiver(syslog *zap.SugaredLogger, fs WriteFS, opts ...RecvOpt) *Receiver {
	o := &RecvOpts{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	return &Receiver{syslog: syslog.Named("file_receiver"), fs: fs, opts: o}
}

func (i *Receiver) State() ReceiveState {
	return i.state
}

func (i *Receiver) Next(req *executorv1.FileTransfer) (err error) {
	defer func() {
		if err != nil {
			if i.fh != nil {
				i.fh.Close()
				i.fh = nil
			}
			i.state = ReceiveStateDone
		}
	}()
	switch i.state {
	case ReceiveStateAwaitingHeaders:
		if req.Header == nil {
			return fmt.Errorf("error header expected")
		}
		i.header = req.Header
		if req.Header.IsDir {
			mode := os.FileMode(req.Header.Mode)
			err = i.fs.MkdirAll(req.Header.DestPath, mode.Perm())
			if err != nil {
				return fmt.Errorf("error making directory: %w", err)
			}
			i.state = ReceiveStateAwaitingTrailer
			if req.Trailer != nil {
				return i.Next(req)
			}
		} else {
			mode := os.FileMode(req.Header.Mode)
			err = i.fs.MkdirAll(filepath.Dir(req.Header.DestPath), 0777)
			if err != nil {
				return fmt.Errorf("error making directory: %w", err)
			}
			file, err := i.fs.OpenFile(req.Header.DestPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}
			i.fh = file
			i.state = ReceiveStateAwaitingBody
			if req.Body != nil {
				return i.Next(req)
			}
		}
	case ReceiveStateAwaitingBody:
		if req.Body == nil {
			i.state = ReceiveStateAwaitingTrailer
			if req.Trailer != nil {
				return i.Next(req)
			}
		} else {
			if i.fh == nil {
				panic("fh should be open")
			}
			_, err = i.fh.Seek(int64(req.Body.Offset), 0)
			if err != nil {
				return fmt.Errorf("error seeking to offset %d: %w", req.Body.Offset, err)
			}
			_, err = i.fh.Write(req.Body.Data)
			if err != nil {
				return fmt.Errorf("error writing data: %w", err)
			}
			if req.Trailer != nil {
				i.state = ReceiveStateAwaitingTrailer
				if req.Trailer != nil {
					return i.Next(req)
				}
			}
		}
	case ReceiveStateAwaitingTrailer:
		if req.Trailer == nil {
			return fmt.Errorf("error trailer expected")
		}
		if req.Trailer.Md5 != nil {
			i.syslog.Warn("MD5 set but verification not implemented")
		}
		i.state = ReceiveStateDone
		return i.Next(req)
	case ReceiveStateDone:
		if i.fh != nil {
			err = i.fh.Close()
			if err != nil {
				return fmt.Errorf("error closing file: %w", err)
			}
			i.fh = nil
		}
		if i.header.IsDir {
			i.syslog.Infow("Received directory", "path", i.header.DestPath)
		} else {
			i.syslog.Infow("Received file", "path", i.header.DestPath, "size", i.header.Size)
		}
		if i.opts.cb != nil {
			i.opts.cb(i.header)
		}
	}
	return nil
}
