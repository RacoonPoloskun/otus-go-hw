package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrInvalidNumberWorkers = errors.New("invalid number workers")
	ErrInvalidMaxErrorCount = errors.New("invalid max error count")
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
)

type Task func() error

type ErrorCount struct {
	value int
	mu    sync.RWMutex
}

func (err *ErrorCount) Add() {
	err.mu.Lock()
	err.value++
	err.mu.Unlock()
}

func (err *ErrorCount) Get() int {
	err.mu.RLock()
	defer err.mu.RUnlock()

	return err.value
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidNumberWorkers
	}

	if m < 0 {
		return ErrInvalidMaxErrorCount
	}

	workers := make(chan func() error, n)
	waitGroup := sync.WaitGroup{}

	errCount := ErrorCount{}

	defer func() {
		close(workers)
		waitGroup.Wait()
	}()

	for _, task := range tasks {
		if errCount.Get() >= m {
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

func runTask(task Task, workers chan func() error, errCount *ErrorCount, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	if task() != nil {
		errCount.Add()
	}

	<-workers
}
