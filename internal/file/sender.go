package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

const chunkSize = 512 * 1024

type SendCallback func(header *executorv1.FileTransferHeader)

type SendOpt interface {
	Apply(*SendOpts)
}

type SendOpts struct {
	cb SendCallback
}

type withSendCallback struct {
	cb SendCallback
}

func (o *withSendCallback) Apply(opts *SendOpts) {
	opts.cb = o.cb
}

func WithSendCallback(cb SendCallback) SendOpt {
	return &withSendCallback{cb: cb}
}

type SendTransport interface {
	Send(*executorv1.FileTransfer) error
}

type SendResult struct{}

type Sender struct {
	syslog    *zap.SugaredLogger
	opts      *SendOpts
	fs        fs.FS
	runtimeID string
}

func NewSender(syslog *zap.SugaredLogger, fs fs.FS, runtimeID string, opts ...SendOpt) *Sender {
	o := &SendOpts{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	return &Sender{syslog: syslog.Named("file_sender"), fs: fs, runtimeID: runtimeID, opts: o}
}

func (f *Sender) Send(stream SendTransport, src string, dest string) (*SendResult, error) {
	if filepath.IsAbs(src) {
		return nil, fmt.Errorf("error src dir must be relative")
	}
	if filepath.IsAbs(dest) {
		return nil, fmt.Errorf("error dest dir must be relative")
	}
	importID := uuid.New().String()

	var isFile, isDir, isDirContents bool
	isGlob, err := isGlob(src)
	if err != nil {
		return nil, fmt.Errorf("error detecting glob in %q: %w", src, err)
	}

	if !isGlob {
		info, err := fs.Stat(f.fs, strings.TrimSuffix(src, "/"))
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("error src %s does not exist", src)
			} else {
				return nil, fmt.Errorf("error stating src: %w", err)
			}
		}
		if strings.HasSuffix(src, "/") || src == "." {
			isDirContents = true
		} else if info.IsDir() {
			isDir = true
		} else {
			isFile = true
		}
	}

	// If src is a single file
	//  Copy to dest, if dest is a directory, then copy to dest/base(src)
	if isFile {
		var finalDest string
		if dest == "" {
			finalDest = src
		} else if strings.HasSuffix(dest, "/") {
			finalDest = filepath.Join(dest, filepath.Base(src))
		} else {
			finalDest = dest
		}
		return &SendResult{}, f.sendFile(stream, importID, src, finalDest)
	}

	// If src is a single directory then
	//	Copy recursively to dest, where dir(src) is substituted for dest in the final dest paths
	if isDir {
		return &SendResult{}, fs.WalkDir(f.fs, src, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			var finalDest string
			if dest == "" {
				finalDest = path
			} else if strings.HasSuffix(dest, "/") {
				finalDest = filepath.Join(dest, strings.TrimPrefix(path, filepath.Dir(src)))
			} else {
				finalDest = filepath.Join(dest, strings.TrimPrefix(path, src))
			}
			if d.IsDir() {
				return f.sendDirectory(stream, importID, path, finalDest)
			} else {
				return f.sendFile(stream, importID, path, finalDest)
			}
		})
	}

	// If src is directory contents e.g. dir/ then
	//  Recurse through everything in the directory and copy to dest/ (dest must be a folder).
	if isDirContents {
		return &SendResult{}, fs.WalkDir(f.fs, strings.TrimSuffix(src, "/"), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// Passing 'foo/' to WalkDir is an error, so we trim the suffix above.
			// Here we make sure we subsequently skip the 'foo' directory.
			if strings.HasSuffix(src, "/") && strings.TrimSuffix(src, "/") == path {
				return nil
			}
			var finalDest string
			if dest == "" {
				finalDest = path
			} else {
				finalDest = filepath.Join(dest, strings.TrimPrefix(path, src))
			}
			if finalDest == "." {
				return nil
			}
			if d.IsDir() {
				return f.sendDirectory(stream, importID, path, finalDest)
			} else {
				return f.sendFile(stream, importID, path, finalDest)
			}
		})
	}

	// If src is a glob
	//  Recurse through every matched file and directory and copy to dest/ (dest must be a folder).
	if isGlob {
		matches, err := doublestar.Glob(f.fs, src, doublestar.WithFailOnIOErrors())
		if err != nil {
			return nil, fmt.Errorf("error expanding glob: %w", err)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("error no matches")
		}
		for _, match := range matches {
			matchDir := filepath.Dir(match)
			err = fs.WalkDir(f.fs, strings.TrimSuffix(match, "/"), func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				var finalDest string
				if dest == "" {
					finalDest = path
				} else {
					if matchDir == "." {
						finalDest = filepath.Join(dest, path)
					} else {
						finalDest = filepath.Join(dest, strings.TrimPrefix(path, matchDir))
					}
				}
				if finalDest == "." {
					return nil
				}
				if d.IsDir() {
					return f.sendDirectory(stream, importID, path, finalDest)
				} else {
					return f.sendFile(stream, importID, path, finalDest)
				}
			})
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return nil, fmt.Errorf("error unable to parse src %q", src)
}

func (f *Sender) sendDirectory(stream SendTransport, importID string, src string, dest string) error {
	fh, err := f.fs.Open(src)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", src, err)
	}
	defer fh.Close()
	info, err := fh.Stat()
	if err != nil {
		return fmt.Errorf("error stating file %s: %w", src, err)
	}
	fileID := uuid.New().String()
	req := &executorv1.FileTransfer{
		RuntimeId: f.runtimeID,
		ImportId:  importID,
		FileId:    fileID,
		Header: &executorv1.FileTransferHeader{
			IsDir:    true,
			SrcPath:  src,
			DestPath: dest,
			Mode:     uint32(info.Mode()),
			Size:     0,
		},
	}
	err = stream.Send(req)
	if err != nil {
		return fmt.Errorf("error sending directory: %w", err)
	}
	f.syslog.Infow("Sent directory", "src", src, "dest", dest, "mode", info.Mode())
	if f.opts.cb != nil {
		f.opts.cb(req.Header)
	}
	return nil
}

func (f *Sender) sendFile(stream SendTransport, importID string, src string, dest string) error {
	fh, err := f.fs.Open(src)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", src, err)
	}
	defer fh.Close()
	info, err := fh.Stat()
	if err != nil {
		return fmt.Errorf("error stating file %s: %w", src, err)
	}
	fileID := uuid.New().String()
	header := &executorv1.FileTransferHeader{
		IsDir:    false,
		SrcPath:  src,
		DestPath: dest,
		Mode:     uint32(info.Mode()),
		Size:     uint64(info.Size()),
	}
	if info.Size() == 0 {
		req := &executorv1.FileTransfer{RuntimeId: f.runtimeID, ImportId: importID, FileId: fileID}
		req.Header = header
		req.Trailer = &executorv1.FileTransferTrailer{}
		err = stream.Send(req)
		if err != nil {
			return fmt.Errorf("error sending file: %w", err)
		}
	} else {
		parts := 0
		for offset := int64(0); offset < info.Size(); {
			buf := make([]byte, chunkSize)
			n, err := fh.Read(buf)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}
			req := &executorv1.FileTransfer{RuntimeId: f.runtimeID, ImportId: importID, FileId: fileID}
			if offset == 0 {
				req.Header = header
			}
			if n > 0 {
				req.Body = &executorv1.FileTransferBody{Offset: uint64(offset), Data: buf[:n]}
			}
			partOffset := offset
			offset += int64(n)
			if offset == info.Size() {
				req.Trailer = &executorv1.FileTransferTrailer{Md5: nil} // TODO
			}
			err = stream.Send(req)
			if err != nil {
				return fmt.Errorf("error sending file: %w", err)
			}
			parts++
			f.syslog.Debugw("Sent file part", "src", src, "part_offset", partOffset, "part_size", n)
		}
	}
	f.syslog.Infow("Sent file", "src", src, "dest", dest, "mode", info.Mode(), "size", info.Size())
	if f.opts.cb != nil {
		f.opts.cb(header)
	}
	return nil
}
