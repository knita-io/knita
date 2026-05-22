package _import

import (
	directorv1 "github.com/knita-io/knita/api/director/v1"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
)

// Opt configures directorv1.ImportOpts.
type Opt func(*directorv1.ImportOpts)

// WithDest sets the destination path for imported files and directories.
func WithDest(destPath string) Opt {
	return func(o *directorv1.ImportOpts) {
		o.DestPath = destPath
	}
}

// WithExcludes excludes files and directories from the import.
func WithExcludes(excludes ...string) Opt {
	return func(o *directorv1.ImportOpts) {
		o.Excludes = append(o.Excludes, excludes...)
	}
}

// WithDisplayName sets the display name for the import.
func WithDisplayName(displayName string) Opt {
	return func(o *directorv1.ImportOpts) {
		o.DisplayName = displayName
	}
}

// WithLabel sets a single label.
func WithLabel(key, value string) Opt {
	return WithLabels(key, value)
}

// WithLabels sets multiple labels from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithLabels(kv ...string) Opt {
	m := runtime.KVMap("WithLabels", kv)
	return func(o *directorv1.ImportOpts) {
		o.Meta = runtime.MergeLabels(o.Meta, m)
	}
}

// WithAnnotation sets a single annotation.
func WithAnnotation(key, value string) Opt {
	return WithAnnotations(key, value)
}

// WithAnnotations sets multiple annotations from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithAnnotations(kv ...string) Opt {
	m := runtime.KVMap("WithAnnotations", kv)
	return func(o *directorv1.ImportOpts) {
		o.Meta = runtime.MergeAnnotations(o.Meta, m)
	}
}
