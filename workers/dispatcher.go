package workers

import (
	"io"
	"runtime"
	"sync"
)

// Slightly modified version of http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

type Doer interface {
	Do()
}

type Dispatcher struct {
	io.Closer

	JobQueue     chan Doer
	CountWorkers int

	workerPool  chan chan Doer
	workersList []*worker
	once        sync.Once
	mu          sync.RWMutex
	quit        bool
}

func NewDispatcher(countWorkers int) *Dispatcher {
	return &Dispatcher{
		CountWorkers: countWorkers,
	}
}

func (d *Dispatcher) Run() {
	d.internalInit()

	d.mu.Lock()
	d.quit = false
	d.mu.Unlock()

	go func() {

		for i := 0; i < len(d.workersList); i++ {
			d.workersList[i].start()
		}

		for {
			d.mu.RLock()
			quit := d.quit
			d.mu.RUnlock()

			if quit {
				for i := 0; i < len(d.workersList); i++ {
					d.workersList[i].stop()
				}
				break
			}

			jobChannel := <-d.workerPool
			jobChannel <- (<-d.JobQueue)
		}
	}()
}

func (d *Dispatcher) Close() error {
	d.internalInit()

	d.mu.Lock()
	d.quit = true
	d.mu.Unlock()

	return nil
}

func (d *Dispatcher) internalInit() {
	d.once.Do(func() {
		if d.CountWorkers == 0 {
			d.CountWorkers = runtime.NumCPU()
		}

		workerPool := make(chan chan Doer, d.CountWorkers)
		workersList := make([]*worker, d.CountWorkers, d.CountWorkers)
		for i := 0; i < d.CountWorkers; i++ {
			workersList[i] = newWorker(i, workerPool)
		}

		d.workerPool = workerPool
		d.workersList = workersList
		d.JobQueue = make(chan Doer, d.CountWorkers)
	})
}
