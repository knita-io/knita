package file

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type WriteFS interface {
	ReadFS() fs.FS
	Directory() string
	MkdirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
}

type File interface {
	io.ReadWriteCloser
	io.Seeker
}

func WriteDirFS(baseDir string) WriteFS {
	return &writeDirFs{baseDir: baseDir}
}

type writeDirFs struct{ baseDir string }

func (r *writeDirFs) ReadFS() fs.FS {
	return os.DirFS(r.baseDir)
}

func (r *writeDirFs) Directory() string {
	return r.baseDir
}

func (r *writeDirFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(filepath.Join(r.baseDir, path), perm)
}

func (r *writeDirFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(filepath.Join(r.baseDir, name), flag, perm)
}

type Logger interface {
	Printf(format string, args ...interface{})
}
