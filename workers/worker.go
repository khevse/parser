package workers

import "sync"

type worker struct {
	id         int
	workerPool chan chan Doer
	jobChannel chan Doer
	quit       bool
	mu         sync.RWMutex
}

func newWorker(id int, workerPool chan chan Doer) *worker {
	return &worker{
		id:         id,
		workerPool: workerPool,
		jobChannel: make(chan Doer),
	}
}

func (w *worker) start() {
	w.mu.Lock()
	w.quit = false
	w.mu.Unlock()

	go func() {
		for {
			w.mu.RLock()
			quit := w.quit
			w.mu.RUnlock()

			if quit {
				return
			}

			w.workerPool <- w.jobChannel
			job := <-w.jobChannel
			job.Do()
		}
	}()
}

func (w *worker) stop() {
	w.mu.Lock()
	w.quit = true
	w.mu.Unlock()
}
