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

type SkipCallback func(path string, isDir bool, excludedBy string)

type SendOpt interface {
	Apply(*SendOpts)
}

type SendOpts struct {
	sendCallback SendCallback
	skipCallback SkipCallback
	excludes     []string
	dest         string
}

type withSendCallback struct {
	cb SendCallback
}

func (o *withSendCallback) Apply(opts *SendOpts) {
	opts.sendCallback = o.cb
}

func WithSendCallback(cb SendCallback) SendOpt {
	return &withSendCallback{cb: cb}
}

type withSkipCallback struct {
	cb SkipCallback
}

func (o *withSkipCallback) Apply(opts *SendOpts) {
	opts.skipCallback = o.cb
}

func WithSkipCallback(cb SkipCallback) SendOpt {
	return &withSkipCallback{cb: cb}
}

type withExcludes struct {
	excludes []string
}

func (o *withExcludes) Apply(opts *SendOpts) {
	opts.excludes = append(opts.excludes, o.excludes...)
}

func WithExcludes(excludes []string) SendOpt {
	return &withExcludes{excludes: excludes}
}

type withDest struct {
	dest string
}

func (o *withDest) Apply(opts *SendOpts) {
	opts.dest = o.dest
}

func WithDest(dest string) SendOpt {
	return &withDest{dest: dest}
}

type SendTransport interface {
	Send(*executorv1.FileTransfer) error
}

type SendResult struct{}

type Sender struct {
	syslog     *zap.SugaredLogger
	opts       *SendOpts
	fs         fs.FS
	transport  SendTransport
	runtimeID  string
	transferID string
}

func NewSender(syslog *zap.SugaredLogger, fs fs.FS, transport SendTransport, runtimeID string, transferID string, opts ...SendOpt) *Sender {
	o := &SendOpts{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	return &Sender{
		syslog:     syslog.Named("file_sender"),
		fs:         fs,
		transport:  transport,
		runtimeID:  runtimeID,
		transferID: transferID,
		opts:       o,
	}
}

func (s *Sender) Send(src string) (*SendResult, error) {
	dest := s.opts.dest
	if filepath.IsAbs(src) {
		return nil, fmt.Errorf("error src must be relative")
	}
	if filepath.IsAbs(dest) {
		return nil, fmt.Errorf("error dest must be relative")
	}
	for _, exclude := range s.opts.excludes {
		if !doublestar.ValidatePattern(exclude) {
			return nil, fmt.Errorf("error invalid exclude pattern: %s", exclude)
		}
	}

	var isFile, isDir, isDirContents bool
	isGlob, err := isGlob(src)
	if err != nil {
		return nil, fmt.Errorf("error detecting glob in %q: %w", src, err)
	}

	if !isGlob {
		info, err := fs.Stat(s.fs, strings.TrimSuffix(src, "/"))
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
		return &SendResult{}, s.filteredSend(false, src, finalDest)
	}

	// If src is a single directory then
	//	Copy recursively to dest, where dir(src) is substituted for dest in the final dest paths
	if isDir {
		return &SendResult{}, fs.WalkDir(s.fs, src, func(path string, d fs.DirEntry, err error) error {
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
			return s.filteredSend(d.IsDir(), path, finalDest)
		})
	}

	// If src is directory contents e.g. dir/ then
	//  Recurse through everything in the directory and copy to dest/ (dest must be a folder).
	if isDirContents {
		return &SendResult{}, fs.WalkDir(s.fs, strings.TrimSuffix(src, "/"), func(path string, d fs.DirEntry, err error) error {
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
			return s.filteredSend(d.IsDir(), path, finalDest)
		})
	}

	// If src is a glob
	//  Recurse through every matched file and directory and copy to dest/ (dest must be a folder).
	if isGlob {
		matches, err := doublestar.Glob(s.fs, src, doublestar.WithFailOnIOErrors())
		if err != nil {
			return nil, fmt.Errorf("error expanding glob: %w", err)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("error no matches")
		}
		for _, match := range matches {
			matchDir := filepath.Dir(match)
			err = fs.WalkDir(s.fs, strings.TrimSuffix(match, "/"), func(path string, d fs.DirEntry, err error) error {
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
				return s.filteredSend(d.IsDir(), path, finalDest)
			})
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return nil, fmt.Errorf("error unable to parse src %q", src)
}

func (s *Sender) filteredSend(isDir bool, src string, dest string) error {
	isExcluded, excludedBy := s.isExcluded(src)
	if isExcluded {
		if s.opts.skipCallback != nil {
			s.opts.skipCallback(src, isDir, excludedBy)
		}
		if isDir {
			s.syslog.Infow("Skipped directory", "path", src, "excluded_by", excludedBy)
			return fs.SkipDir
		}
		s.syslog.Infow("Skipped file", "path", src, "excluded_by", excludedBy)
		return nil
	}
	if isDir {
		return s.sendDirectory(src, dest)
	} else {
		return s.sendFile(src, dest)
	}
}

func (s *Sender) isExcluded(path string) (bool, string) {
	for _, exclude := range s.opts.excludes {
		isGlob, _ := isGlob(exclude)
		if isGlob {
			match, _ := doublestar.Match(exclude, path)
			if match {
				return true, exclude
			}
		} else if path == exclude {
			return true, exclude
		} else {
			up := "../"
			rel, err := filepath.Rel(exclude, path)
			if err != nil {
				return false, ""
			}
			if !strings.HasPrefix(rel, up) && rel != ".." {
				return true, exclude
			}
		}
	}
	return false, ""
}

func (s *Sender) sendDirectory(src string, dest string) error {
	fh, err := s.fs.Open(src)
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
		RuntimeId:  s.runtimeID,
		TransferId: s.transferID,
		FileId:     fileID,
		Header: &executorv1.FileTransferHeader{
			IsDir:    true,
			SrcPath:  src,
			DestPath: dest,
			Mode:     uint32(info.Mode()),
			Size:     0,
		},
	}
	err = s.transport.Send(req)
	if err != nil {
		return fmt.Errorf("error sending directory: %w", err)
	}
	s.syslog.Infow("Sent directory", "src", src, "dest", dest, "mode", info.Mode())
	if s.opts.sendCallback != nil {
		s.opts.sendCallback(req.Header)
	}
	return nil
}

func (s *Sender) sendFile(src string, dest string) error {
	fh, err := s.fs.Open(src)
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
		req := &executorv1.FileTransfer{RuntimeId: s.runtimeID, TransferId: s.transferID, FileId: fileID}
		req.Header = header
		req.Trailer = &executorv1.FileTransferTrailer{}
		err = s.transport.Send(req)
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
			req := &executorv1.FileTransfer{RuntimeId: s.runtimeID, TransferId: s.transferID, FileId: fileID}
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
			err = s.transport.Send(req)
			if err != nil {
				return fmt.Errorf("error sending file: %w", err)
			}
			parts++
			s.syslog.Debugw("Sent file part", "src", src, "part_offset", partOffset, "part_size", n)
		}
	}
	s.syslog.Infow("Sent file", "src", src, "dest", dest, "mode", info.Mode(), "size", info.Size())
	if s.opts.sendCallback != nil {
		s.opts.sendCallback(header)
	}
	return nil
}
