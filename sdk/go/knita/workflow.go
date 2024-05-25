package knita

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

var (
	errorType     = reflect.TypeOf((*error)(nil)).Elem()
	errQueueEmpty = errors.New("error queue is empty")
)

type ioID string
type jobID string

type jobDescriptor struct {
	id       jobID
	fn       interface{}
	fnT      reflect.Type
	fnV      reflect.Value
	needs    []ioID
	provides []ioID
}

type Workflow struct {
	log                Log
	fatalFunc          FatalFunc
	director           Director
	mu                 sync.Mutex
	jobSeqNo           int
	typeToIOID         map[reflect.Type]ioID
	resolvedInputsByID map[ioID]interface{}
	providedInputsByID map[ioID]struct{}
	jobByID            map[jobID]jobDescriptor
}

type job struct {
	fn  interface{}
	in  interface{}
	out interface{}
}

func Job[In any, Out any](fn func(input In) (Out, error)) job {
	var in In
	var out Out
	return job{fn: fn, in: in, out: out}
}

func (w *Workflow) WithInput(input interface{}) *Workflow {
	w.mu.Lock()
	defer w.mu.Unlock()
	id := w.recordIO(reflect.TypeOf(input))
	_, ok := w.providedInputsByID[id]
	if ok {
		w.fatalFunc(fmt.Errorf("same input registered twice"))
	}
	w.providedInputsByID[id] = struct{}{}
	w.resolvedInputsByID[id] = input
	w.director.ResolveInput(id)
	return w
}

func (w *Workflow) WithJob(f job) *Workflow {
	w.mu.Lock()
	defer w.mu.Unlock()
	needs, err := w.recordRootIO(f.in)
	if err != nil {
		w.fatalFunc(err)
	}
	provides, err := w.recordRootIO(f.out)
	if err != nil {
		w.fatalFunc(err)
	}
	for _, id := range provides {
		_, ok := w.providedInputsByID[id]
		if ok {
			w.fatalFunc(fmt.Errorf("same input registered twice"))
		}
		w.providedInputsByID[id] = struct{}{}
	}
	w.jobSeqNo++
	job := jobDescriptor{
		id:       jobID(strconv.FormatInt(int64(w.jobSeqNo), 10)),
		fn:       f.fn,
		fnT:      reflect.TypeOf(f.fn),
		fnV:      reflect.ValueOf(f.fn),
		needs:    needs,
		provides: provides,
	}
	w.jobByID[job.id] = job
	w.director.Enqueue(job.id, job.needs, job.provides)
	w.log.Printf("workflow: Added job: id=%s, needs=%v, provides=%v\n", job.id, job.needs, job.provides)
	return w
}

func (w *Workflow) MustRun() {
	err := w.Run()
	if err != nil {
		w.fatalFunc(fmt.Errorf("error in workflow: %v", err))
	}
}

func (w *Workflow) Run() error {
	w.log.Printf("workflow: running\n")
	for {
		jobID, err := w.director.Dequeue()
		if err != nil {
			if errors.Is(err, errQueueEmpty) {
				w.log.Printf("workflow: finished\n")
				return nil
			}
			// The graph has unreachable nodes in it.
			w.fatalFunc(fmt.Errorf("error dequeuing next job: %w", err))
		}
		w.mu.Lock()
		job, ok := w.jobByID[jobID]
		w.mu.Unlock()
		if !ok {
			// The SDK and the director are out of sync.
			w.fatalFunc(fmt.Errorf("error locating dequeued job: %s", jobID))
		}
		err = w.validateJobFn(job)
		if err != nil {
			// WithJob allowed an invalid job signature through (even though we're using
			// generics and the type system should have prevented it).
			w.fatalFunc(fmt.Errorf("error validating job function: %w", err))
		}
		go func(job jobDescriptor) {
			w.log.Printf("workflow: Running job: id=%s\n", job.id)
			args := []reflect.Value{w.resolveJobInput(job.fnT.In(0))}
			results := job.fnV.Call(args)
			if results[1].Interface() != nil {
				err := results[1].Interface().(error)
				w.fatalFunc(fmt.Errorf("job failed: %v", err)) // TODO return error
			}
			out := results[0].Interface()
			_, err = w.recordRootInput(out)
			if err != nil {
				// The graph provided the same input multiple times, making it ambiguous.
				w.fatalFunc(fmt.Errorf("error recording input: %v", err))
			}
			w.director.Complete(job.id)
		}(job)
	}
}

func (w *Workflow) validateJobFn(job jobDescriptor) error {
	if job.fnT.NumIn() != 1 {
		return fmt.Errorf("invalid number of arguments to job function")
	}
	in1 := job.fnT.In(0)
	if !(in1.Kind() == reflect.Struct || (in1.Kind() == reflect.Ptr && in1.Elem().Kind() == reflect.Struct)) {
		return fmt.Errorf("expected first arg of job function to be a struct or a pointer to a struct")
	}
	if job.fnT.NumOut() != 2 {
		return fmt.Errorf("invalid number of return values from job function")
	}
	out1 := job.fnT.Out(0)
	if !(out1.Kind() == reflect.Struct || (out1.Kind() == reflect.Ptr && out1.Elem().Kind() == reflect.Struct)) {
		return fmt.Errorf("expected first return value of job function to be a struct or a pointer to a struct")
	}
	out2 := job.fnT.Out(1)
	if !out2.Implements(errorType) {
		return fmt.Errorf("expected second return value of job function to be an error")
	}
	return nil
}

func (w *Workflow) resolveJobInput(inputT reflect.Type) reflect.Value {
	var ptr bool
	if inputT.Kind() == reflect.Ptr {
		ptr = true
		inputT = inputT.Elem()
	}
	if inputT.Kind() != reflect.Struct {
		panic("expected struct or pointer to a struct")
	}
	v := reflect.New(inputT)
	for i := 0; i < inputT.NumField(); i++ {
		inputID := w.typeToIOID[inputT.Field(i).Type]
		input := w.resolvedInputsByID[inputID]
		v.Elem().Field(i).Set(reflect.ValueOf(input))
	}
	if !ptr {
		return v.Elem()
	}
	return v
}

func (w *Workflow) recordRootIO(io interface{}) ([]ioID, error) {
	w.recordIO(reflect.TypeOf(io))
	t, err := w.getStructT(reflect.TypeOf(io))
	if err != nil {
		return nil, err
	}
	var ids []ioID
	for i := 0; i < t.NumField(); i++ {
		id := w.recordIO(t.Field(i).Type)
		ids = append(ids, id)
	}
	return ids, nil
}

func (w *Workflow) recordIO(t reflect.Type) ioID {
	id, ok := w.typeToIOID[t]
	if !ok {
		id = ioID(t.String())
		w.typeToIOID[t] = id
	}
	return id
}

func (w *Workflow) recordRootInput(input interface{}) ([]ioID, error) {
	_, err := w.recordInput(input)
	if err != nil {
		return nil, err
	}
	t, err := w.getStructT(reflect.TypeOf(input))
	if err != nil {
		return nil, err
	}
	ids := []ioID{w.recordIO(t)}
	for i := 0; i < t.NumField(); i++ {
		field := reflect.ValueOf(input).Elem().Field(i).Interface()
		id, err := w.recordInput(field)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (w *Workflow) recordInput(input interface{}) (ioID, error) {
	id := w.recordIO(reflect.TypeOf(input))
	_, ok := w.resolvedInputsByID[id]
	if ok {
		return "", fmt.Errorf("error input %s recorded twice; each input must be a unique type", id)
	}
	w.resolvedInputsByID[id] = input
	w.director.ResolveInput(id)
	w.log.Printf("workflow: Input now available: %s\n", id)
	return id, nil
}

// getStructT checks that t is a struct or pointer to a struct, and in either case returns the struct.
func (w *Workflow) getStructT(t reflect.Type) (reflect.Type, error) {
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
	}
	if !(t.Kind() == reflect.Struct) {
		return nil, fmt.Errorf("expected input to be a struct or pointer to a struct")
	}
	return t, nil
}
