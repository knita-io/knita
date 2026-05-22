package runtime

import (
	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

// MergeLabels merges kv into meta.Labels, allocating meta and its Labels map
// as needed. Returns the (possibly newly allocated) meta so callers can write:
//
//	o.Meta = runtime.MergeLabels(o.Meta, kv)
func MergeLabels(meta *executorv1.OptsMeta, kv map[string]string) *executorv1.OptsMeta {
	if len(kv) == 0 {
		return meta
	}
	if meta == nil {
		meta = &executorv1.OptsMeta{}
	}
	if meta.Labels == nil {
		meta.Labels = make(map[string]string, len(kv))
	}
	for k, v := range kv {
		meta.Labels[k] = v
	}
	return meta
}

// MergeAnnotations is the annotations counterpart of MergeLabels.
func MergeAnnotations(meta *executorv1.OptsMeta, kv map[string]string) *executorv1.OptsMeta {
	if len(kv) == 0 {
		return meta
	}
	if meta == nil {
		meta = &executorv1.OptsMeta{}
	}
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string, len(kv))
	}
	for k, v := range kv {
		meta.Annotations[k] = v
	}
	return meta
}

// KVMap converts an alternating key/value slice into a map. Panics on odd length.
// Exposed so the per-call sub-packages (exec, import, export) can share the same
// argument-parsing semantics as the runtime package's own WithLabels.
func KVMap(name string, kv []string) map[string]string {
	if len(kv)%2 != 0 {
		panic(name + ": expected even number of key/value args")
	}
	m := make(map[string]string, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return m
}
