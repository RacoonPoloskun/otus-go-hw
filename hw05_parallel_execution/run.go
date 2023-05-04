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

	workers := make(chan Task, len(tasks))
	waitGroup := sync.WaitGroup{}

	errCount := ErrorCount{}

	for _, task := range tasks {
		workers <- task
	}
	close(workers)

	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for {
				if task, ok := <-workers; ok && errCount.Get() < m {
					if res := task(); res != nil {
						errCount.Add()
					}
				} else {
					return
				}
			}
		}()
	}

	waitGroup.Wait()

	if errCount.Get() >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
