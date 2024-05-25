package knita

import (
	"fmt"
	"sync"
)

type Director interface {
	ResolveInput(ioID)
	Enqueue(id jobID, needs []ioID, provides []ioID)
	Dequeue() (jobID, error)
	Complete(id jobID)
}

type jobInfo struct {
	id       jobID
	needs    []ioID
	provides []ioID
}

type director struct {
	mu             sync.Mutex
	resolvedInputs map[ioID]struct{}
	jobs           map[jobID]jobInfo
	active         map[jobID]jobInfo
	next           chan jobID
}

func newDirector() *director {
	return &director{
		mu:             sync.Mutex{},
		resolvedInputs: map[ioID]struct{}{},
		jobs:           map[jobID]jobInfo{},
		active:         map[jobID]jobInfo{},
		next:           make(chan jobID, 100000),
	}
}

func (d *director) ResolveInput(id ioID) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.resolvedInputs[id]
	if ok {
		panic("Same input resolved twice")
	}
	d.resolvedInputs[id] = struct{}{}
	d.stokeNext()
}

func (d *director) Enqueue(id jobID, needs []ioID, provides []ioID) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.jobs[id]
	if ok {
		panic("same job enqueued twice")
	}
	d.jobs[id] = jobInfo{
		id:       id,
		needs:    needs,
		provides: provides,
	}
	d.stokeNext()
}

func (d *director) Dequeue() (jobID, error) {
	d.mu.Lock()

	if len(d.active) == 0 && len(d.jobs) > 0 {
		d.mu.Unlock()
		return "", fmt.Errorf("error unreachable jobs")
	}

	if len(d.active) == 0 && len(d.jobs) == 0 {
		d.mu.Unlock()
		return "", nil
	}

	d.mu.Unlock()

	next, ok := <-d.next
	if !ok {
		return "", errQueueEmpty
	}
	return next, nil
}

func (d *director) Complete(jobID jobID) {
	d.mu.Lock()
	defer d.mu.Unlock()
	job, ok := d.active[jobID]
	if !ok {
		panic("job not found")
	}
	delete(d.active, jobID)
	for _, inputID := range job.provides {
		d.resolvedInputs[inputID] = struct{}{}
	}
	d.stokeNext()
	if len(d.active) == 0 && len(d.jobs) == 0 {
		close(d.next)
	}
}

func (d *director) stokeNext() {
	var readyJobs []jobInfo
	for _, job := range d.jobs {
		ready := true
		for _, ioID := range job.needs {
			if _, ok := d.resolvedInputs[ioID]; !ok {
				ready = false
				break
			}
		}
		if ready {
			readyJobs = append(readyJobs, job)
		}
	}
	for _, job := range readyJobs {
		delete(d.jobs, job.id)
		d.active[job.id] = job
		d.next <- job.id
	}
}
