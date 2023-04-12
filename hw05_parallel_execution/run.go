package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrInvalidNumberWorkers = errors.New("invalid number workers")
	ErrInvalidMaxErrorCount = errors.New("invalid max error count")
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidNumberWorkers
	}

	if m < 0 {
		return ErrInvalidMaxErrorCount
	}

	workers := make(chan func() error, n)
	waitGroup := sync.WaitGroup{}

	var errCount int32

	defer func() {
		close(workers)
		waitGroup.Wait()
	}()

	for _, task := range tasks {
		if errCount >= int32(m) {
			return ErrErrorsLimitExceeded
		}

		workers <- func() error {
			return nil
		}

		waitGroup.Add(1)
		go runTask(task, workers, &errCount, &waitGroup)
	}

	return nil
}

func runTask(task Task, workers chan func() error, errCount *int32, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	if task() != nil {
		atomic.AddInt32(errCount, 1)
	}

	<-workers
}
