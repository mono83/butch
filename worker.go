package butch

import (
	"bytes"
	"io"
	"runtime"
	"time"
)

// Worker represents batch worker
// Consumers must supply name of worker and run function
type Worker struct {
	state            int
	start, stop      time.Time
	stdout, stderr   *bytes.Buffer
	lastMsg, lastErr string
	Name             string
	Run              func(pool *Pool, worker Worker, writer io.Writer, errWriter io.Writer)
}

func (w *Worker) forkStart(p *Pool) {
	go func() {
		runtime.LockOSThread()
		w.Start(p)
		runtime.UnlockOSThread()
	}()
}

// Start starts worker
func (w *Worker) Start(p *Pool) {
	w.start = time.Now()
	w.stdout = bytes.NewBuffer(nil)
	w.stderr = bytes.NewBuffer(nil)

	stdout := &workerByteBuffer{
		full:   w.stdout,
		chunk:  bytes.NewBuffer(nil),
		target: &w.lastMsg,
	}
	stderr := &workerByteBuffer{
		full:   w.stderr,
		chunk:  bytes.NewBuffer(nil),
		target: &w.lastErr,
	}

	w.state = 1
	w.Run(p, *w, stdout, stderr)
	stdout.done()
	stderr.done()
	w.state = 2
	w.stop = time.Now()
}

// StartedAt returns worker start time
func (w Worker) StartedAt() time.Time {
	return w.start
}

// DoneAt returns workers done time
func (w Worker) DoneAt() time.Time {
	return w.stop
}

// Elapsed returns worker's elapsed time
func (w Worker) Elapsed() time.Duration {
	if w.IsWaiting() {
		return time.Duration(0)
	} else if w.IsDone() {
		return w.stop.Sub(w.start)
	}

	return time.Now().Sub(w.start)
}

// IsWaiting returns true if worker wasn't started
func (w Worker) IsWaiting() bool {
	return w.state == 0
}

// IsDone return true if worker completed
func (w Worker) IsDone() bool {
	return w.state > 1
}
