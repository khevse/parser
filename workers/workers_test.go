package workers

import (
	"runtime"
	"testing"
)

type JobResult struct {
	Thread int
	Number int
}

type Job struct {
	Thread  int
	Number  int
	Results chan *JobResult
}

func (j *Job) Do() {
	j.Results <- &JobResult{
		Thread: j.Thread,
		Number: j.Number,
	}
}

func TestWorkers(t *testing.T) {

	const (
		COUNT_THREADS = 1000
		COUNT_JOBS    = 1000
	)

	results := make(chan *JobResult, COUNT_JOBS)

	dispatcher := new(Dispatcher)
	dispatcher.Run()
	defer dispatcher.Close()

	if dispatcher.CountWorkers != runtime.NumCPU() {
		t.Error("Fail:", dispatcher.CountWorkers)
		return
	}

	fnWork := func(thread int, queue chan Doer) {
		for i := 0; i < COUNT_JOBS; i++ {
			job := &Job{
				Thread:  thread,
				Number:  i,
				Results: results,
			}

			queue <- job
		}
	}

	for i := 0; i < COUNT_THREADS; i++ {
		go fnWork(i, dispatcher.JobQueue)
	}

	total := make([]int, COUNT_THREADS)

	for r := range results {
		total[r.Thread] += 1

		count := 0
		for _, v := range total {
			count += v
		}

		if count == COUNT_THREADS*COUNT_JOBS {
			for thread, count := range total {
				if count != COUNT_JOBS {
					t.Error("Fail:", thread, count)
				}
			}

			return
		}
	}

	results <- nil
	if r := results; r != nil {
		t.Error(r)
	}
}
