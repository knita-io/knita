package director

import (
	"errors"
	"fmt"
)

var (
	errNoReadyJobs = errors.New("error no ready jobs")
	errQueueEmpty  = errors.New("error queue empty")
)

type job struct {
	id       string
	needs    []string
	provides []string
}

type workflow struct {
	log           *Log
	firstJobAdded bool
	queuedJobs    map[string]job
	activeJobs    map[string]job
	provides      map[string]string
	resolved      map[string]string
}

func newWorkflow(log *Log) *workflow {
	log.Printf("workflow: Started")
	return &workflow{
		log:        log,
		queuedJobs: map[string]job{},
		activeJobs: map[string]job{},
		provides:   map[string]string{},
		resolved:   map[string]string{},
	}
}

func (w *workflow) Enqueue(id string, needs []string, provides []string) error {
	job := job{id: id, needs: needs, provides: provides}
	for _, io := range job.provides {
		otherJobID, ok := w.provides[io]
		if ok {
			return fmt.Errorf("error %s is provided by jobs %s and %s; job outputs must be unique",
				io, otherJobID, job.id)
		}
		w.provides[io] = job.id
	}
	w.queuedJobs[job.id] = job
	w.firstJobAdded = true
	w.log.Printf("workflow: Added job: id=%s, needs=%v, provides=%v", job.id, job.needs, job.provides)
	return nil
}

func (w *workflow) AddInput(source string, io string) error {
	otherSource, ok := w.resolved[io]
	if ok {
		return fmt.Errorf("error %s is provided by %q and %q; inputs must be unique", io, otherSource, source)
	}
	w.resolved[io] = source
	w.log.Printf("workflow: Input now available: %s (from %q)", io, source)
	return nil
}

func (w *workflow) CompleteJob(id string) error {
	job, ok := w.activeJobs[id]
	if !ok {
		return fmt.Errorf("error job %s was not found", id)
	}
	delete(w.activeJobs, id)
	for _, io := range job.provides {
		err := w.AddInput(fmt.Sprintf("CompleteJob(%s)", id), io)
		if err != nil {
			return err
		}
	}
	w.log.Printf("workflow: Completed job: id=%s\n", job.id)
	return nil
}

func (w *workflow) GetDependencies(id string) ([]string, error) {
	job, ok := w.activeJobs[id]
	if !ok {
		return nil, fmt.Errorf("error job %s was not found", id)
	}
	dependenciesM := map[string]struct{}{}
	for _, need := range job.needs {
		providedBy := w.provides[need]
		if providedBy != "" { // Some inputs are provided outside a job
			dependenciesM[providedBy] = struct{}{}
		}
	}
	dependencies := make([]string, 0, len(dependenciesM))
	for k := range dependenciesM {
		dependencies = append(dependencies, k)
	}
	return dependencies, nil
}

func (w *workflow) Dequeue() (string, error) {
	var readyJob *job
	for _, job := range w.queuedJobs {
		ready := true
		for _, io := range job.needs {
			if _, ok := w.resolved[io]; !ok {
				ready = false
				break
			}
		}
		if ready {
			readyJob = &job
			break
		}
	}
	if readyJob != nil {
		delete(w.queuedJobs, readyJob.id)
		w.activeJobs[readyJob.id] = *readyJob
		w.log.Printf("workflow: Running job: id=%s", readyJob.id)
		return readyJob.id, nil
	}
	if len(w.activeJobs) != 0 || !w.firstJobAdded {
		return "", errNoReadyJobs
	}
	if len(w.queuedJobs) != 0 {
		return "", fmt.Errorf("error workflow contains unreachable jobs") // TODO summarize what jobs are waiting on what
	}
	w.log.Printf("workflow: Finished")
	return "", errQueueEmpty
}
