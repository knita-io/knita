package _import

import (
	directorv1 "github.com/knita-io/knita/api/director/v1"
)

type Opt interface {
	Apply(opts *directorv1.ImportOpts)
}

type withDest struct {
	destPath string
}

func (o *withDest) Apply(opts *directorv1.ImportOpts) {
	opts.DestPath = o.destPath
}

// WithDest sets the destination path for imported files and directories.
func WithDest(destPath string) Opt {
	return &withDest{destPath: destPath}
}

type withExcludes struct {
	excludes []string
}

func (o *withExcludes) Apply(opts *directorv1.ImportOpts) {
	opts.Excludes = append(opts.Excludes, o.excludes...)
}

// WithExcludes excludes files and directories from the import.
func WithExcludes(excludes ...string) Opt {
	return &withExcludes{excludes: excludes}
}
