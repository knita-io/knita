package knita

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	directorv1 "github.com/knita-io/knita/api/director/v1"
)

var (
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
	clientType = reflect.TypeOf(&Client{})
)

type IORecorder interface {
	IOData() json.RawMessage
}

type jobID string

type jobDescriptor struct {
	id     jobID
	fn     interface{}
	fnT    reflect.Type
	fnV    reflect.Value
	fnArgT reflect.Type
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

type Workflow struct {
	syslog          Log
	fatalFunc       FatalFunc
	client          directorv1.Director_WorkflowClient
	mu              sync.Mutex
	jobSeqNo        int
	jobsByID        map[jobID]jobDescriptor
	inputValuesByID map[string]interface{}
}

func newWorkflow(syslog Log, fatalFunc FatalFunc, client *Client) (*Workflow, error) {
	workflowClient, err := client.client.Workflow(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error making workflow client: %w", err)
	}
	workflow := &Workflow{
		syslog:          syslog,
		fatalFunc:       fatalFunc,
		client:          workflowClient,
		jobsByID:        map[jobID]jobDescriptor{},
		inputValuesByID: map[string]interface{}{},
	}
	return workflow.WithInput(client).WithInput(workflow), nil
}

func (w *Workflow) WithInput(input interface{}) *Workflow {
	w.mu.Lock()
	defer w.mu.Unlock()
	id := w.getID(reflect.TypeOf(input))
	update := &directorv1.WorkflowUpdate{
		Payload: &directorv1.WorkflowUpdate_AddInput{
			AddInput: &directorv1.WorkflowAddInput{InputId: id},
		},
	}
	w.mustSendUpdate(update)
	w.inputValuesByID[id] = input
	return w
}

func (w *Workflow) WithJob(job job) *Workflow {
	w.mu.Lock()
	defer w.mu.Unlock()
	fnT := reflect.TypeOf(job.fn)
	if err := w.validateJobFn(fnT); err != nil {
		w.fatalFunc(err)
	}
	w.jobSeqNo++
	desc := jobDescriptor{
		id:     jobID(fmt.Sprintf("%d", w.jobSeqNo)),
		fn:     job.fn,
		fnT:    fnT,
		fnV:    reflect.ValueOf(job.fn),
		fnArgT: fnT.In(0),
	}
	w.jobsByID[desc.id] = desc
	update := &directorv1.WorkflowUpdate{
		Payload: &directorv1.WorkflowUpdate_AddJob{
			AddJob: &directorv1.WorkflowAddJob{
				JobId: string(desc.id),
				// A job needs all fields on its input struct to be filled
				Needs: w.getIDsOfFields(job.in),
				// A job provides a) its top level output struct, and b) all the fields contained within that struct
				Provides: append(w.getIDsOfFields(job.out), w.getID(reflect.TypeOf(job.out))),
			},
		},
	}
	w.mustSendUpdate(update)
	return w
}

func (w *Workflow) MustRun() {
	err := w.Run()
	if err != nil {
		w.fatalFunc(fmt.Errorf("error in workflow: %v", err))
	}
}

func (w *Workflow) Run() error {
	for {
		signal, err := w.client.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			w.fatalFunc(err)
		}
		next, ok := signal.Payload.(*directorv1.WorkflowSignal_JobReady)
		if !ok {
			w.syslog.Printf("Ignoring unknown signal from workflow director: %s\n", signal.Payload)
			continue
		}
		w.mu.Lock()
		job, ok := w.jobsByID[jobID(next.JobReady.JobId)]
		w.mu.Unlock()
		if !ok {
			w.fatalFunc(fmt.Errorf("error locating dequeued job: %s", next.JobReady.JobId))
		}
		go func(job jobDescriptor) {
			input := w.makeJobInput(job)
			w.updateJobStart(job, input)
			args := []reflect.Value{input}
			start := time.Now()
			results := job.fnV.Call(args)
			duration := time.Now().Sub(start)
			if results[1].Interface() != nil {
				err := results[1].Interface().(error)
				w.fatalFunc(fmt.Errorf("job failed: %v", err)) // TODO return error
			}
			output := results[0].Interface()
			w.mu.Lock()
			w.recordJobOutput(output)
			w.mu.Unlock()
			w.updateJobComplete(job, duration, output)
		}(job)
	}
}

func (w *Workflow) validateJobFn(fnT reflect.Type) error {
	if fnT.NumIn() != 1 {
		return fmt.Errorf("invalid number of arguments to job function")
	}
	in1 := fnT.In(0)
	if !(in1.Kind() == reflect.Struct || (in1.Kind() == reflect.Ptr && in1.Elem().Kind() == reflect.Struct)) {
		return fmt.Errorf("expected first arg of job function to be a struct or a pointer to a struct")
	}
	if fnT.NumOut() != 2 {
		return fmt.Errorf("invalid number of return values from job function")
	}
	out1 := fnT.Out(0)
	if !(out1.Kind() == reflect.Struct || (out1.Kind() == reflect.Ptr && out1.Elem().Kind() == reflect.Struct)) {
		return fmt.Errorf("expected first return value of job function to be a struct or a pointer to a struct")
	}
	out2 := fnT.Out(1)
	if !out2.Implements(errorType) {
		return fmt.Errorf("expected second return value of job function to be an error")
	}
	return nil
}

func (w *Workflow) getID(of reflect.Type) string {
	return of.String()
}

func (w *Workflow) getIDsOfFields(from interface{}) []string {
	var ids []string
	t := reflect.TypeOf(from)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Type
		ids = append(ids, w.getID(field))
	}
	return ids
}

func (w *Workflow) makeJobInput(job jobDescriptor) reflect.Value {
	fnArgT := job.fnArgT
	if fnArgT.Kind() == reflect.Ptr {
		fnArgT = fnArgT.Elem()
	}
	value := reflect.New(fnArgT)
	for i := 0; i < fnArgT.NumField(); i++ {
		id := w.getID(fnArgT.Field(i).Type)
		fieldValue := w.inputValuesByID[id]
		// Special case - if the input is a knita.Client then inject the job id into
		// the client before passing it to the job. This is how we tag client interactions
		// with the correct job id.
		if fnArgT.Field(i).Type == clientType {
			client := fieldValue.(*Client).clone()
			client.jobID = string(job.id)
			fieldValue = client
		}
		value.Elem().Field(i).Set(reflect.ValueOf(fieldValue))
	}
	if job.fnArgT.Kind() != reflect.Ptr {
		return value.Elem()
	}
	return value
}

func (w *Workflow) recordJobOutput(out interface{}) {
	t := reflect.TypeOf(out)
	// Record the top level output struct value
	id := w.getID(t)
	w.inputValuesByID[id] = out
	// Record the values of the fields of the top level output struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(out)
	for i := 0; i < t.NumField(); i++ {
		field := v.Elem().Field(i)
		id := w.getID(field.Type())
		w.inputValuesByID[id] = field.Interface()
	}
}

func (w *Workflow) updateJobStart(job jobDescriptor, input interface{}) {
	var inputData json.RawMessage
	recorder, ok := input.(IORecorder)
	if ok {
		inputData = recorder.IOData()
	}
	update := &directorv1.WorkflowUpdate{
		Payload: &directorv1.WorkflowUpdate_StartJob{
			StartJob: &directorv1.WorkflowStartJob{
				JobId:     string(job.id),
				InputData: inputData,
			},
		},
	}
	w.mustSendUpdate(update)
}

func (w *Workflow) updateJobComplete(job jobDescriptor, duration time.Duration, output interface{}) {
	var outputData json.RawMessage
	recorder, ok := output.(IORecorder)
	if ok {
		outputData = recorder.IOData()
	}
	update := &directorv1.WorkflowUpdate{
		Payload: &directorv1.WorkflowUpdate_CompleteJob{
			CompleteJob: &directorv1.WorkflowCompleteJob{
				JobId:      string(job.id),
				Duration:   durationpb.New(duration),
				OutputData: outputData,
			},
		},
	}
	w.mustSendUpdate(update)
}

func (w *Workflow) mustSendUpdate(update *directorv1.WorkflowUpdate) {
	err := w.client.Send(update)
	if err != nil {
		w.fatalFunc(err)
	}
}
