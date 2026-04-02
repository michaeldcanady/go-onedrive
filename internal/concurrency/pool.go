package concurrency

import (
	"sync"
)

const (
	// DefaultCapacity is the default maximum number of concurrent tasks.
	DefaultCapacity = 5
)

// Option defines a functional configuration for the WorkerPool.
type Option func(*WorkerPool)

// WithCapacity sets the maximum number of concurrent tasks for the pool.
func WithCapacity(capacity int) Option {
	return func(wp *WorkerPool) {
		if capacity > 0 {
			wp.capacity = capacity
		}
	}
}

// WorkerPool manages a pool of goroutines to execute tasks concurrently.
type WorkerPool struct {
	capacity int
	sem      chan struct{}
	wg       sync.WaitGroup
	once     sync.Once
}

// NewWorkerPool creates a new WorkerPool with the provided options.
func NewWorkerPool(opts ...Option) *WorkerPool {
	wp := &WorkerPool{
		capacity: DefaultCapacity,
	}

	for _, opt := range opts {
		opt(wp)
	}

	return wp
}

// init ensures the semaphore is initialized with the correct capacity.
func (wp *WorkerPool) init() {
	wp.once.Do(func() {
		wp.sem = make(chan struct{}, wp.capacity)
	})
}

// Submit adds a task to the worker pool for execution.
// It blocks if the pool is at capacity until a worker is available.
func (wp *WorkerPool) Submit(task func()) {
	wp.init()
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
