package butch

import (
	"runtime"
	"sync"
	"time"
)

// Pool is executions pool
type Pool struct {
	m         sync.Mutex
	startedAt time.Time
	workers   []*Worker
}

// AddWorker registers new worker
func (p *Pool) AddWorker(w *Worker) {
	p.m.Lock()
	defer p.m.Unlock()
	if w != nil {
		p.workers = append(p.workers, w)
	}
}

// Workers returns slice of worker structs, not pointers
func (p *Pool) Workers() []Worker {
	p.m.Lock()
	defer p.m.Unlock()
	response := make([]Worker, len(p.workers))
	for i, w := range p.workers {
		response[i] = *w
	}

	return response
}

// Elapsed returns pool execution duration
func (p *Pool) Elapsed() time.Duration {
	return time.Now().Sub(p.startedAt)
}

// Start starts pool
func (p *Pool) Start() {
	p.startedAt = time.Now()

	// Building draw area
	da := &ProcessRenderer{pool: p}
	running := true
	da.Paint()

	go func() {
		for running {
			da.Paint()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	runtime.GOMAXPROCS(len(p.workers) + 10)

	for {
		allDone := true
		p.m.Lock()
		for _, w := range p.workers {
			if w.IsWaiting() {
				w.forkStart(p)
				allDone = false
			} else if !w.IsDone() {
				allDone = false
			}
		}
		p.m.Unlock()

		if allDone {
			break
		}
	}

	da.Paint()
}
