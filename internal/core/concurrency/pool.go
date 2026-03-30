package concurrency

import (
	"sync"
)

// WorkerPool manages a pool of goroutines to execute tasks concurrently.
type WorkerPool struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// NewWorkerPool creates a new WorkerPool with a given capacity.
func NewWorkerPool(capacity int) *WorkerPool {
	if capacity <= 0 {
		capacity = 1 // Default to at least 1 worker if capacity is invalid
	}
	return &WorkerPool{
		sem: make(chan struct{}, capacity),
	}
}

// Submit adds a task to the worker pool for execution.
// It blocks if the pool is at capacity until a worker is available.
func (wp *WorkerPool) Submit(task func()) {
	wp.sem <- struct{}{} // Acquire a worker slot
	wp.wg.Add(1)
	go func() {
		defer func() {
			<-wp.sem // Release the worker slot
			wp.wg.Done()
		}()
		task()
	}()
}

// Wait blocks until all submitted tasks have completed.
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
