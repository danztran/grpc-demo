package util

import "sync"

type WorkerPool struct {
	wg   sync.WaitGroup
	pool chan func()
}

// NewWorkerPool return a new worker pool
func NewWorkerPool(size int) *WorkerPool {
	p := &WorkerPool{
		wg:   sync.WaitGroup{},
		pool: make(chan func()),
	}

	p.startPool(size)

	return p
}

// Run send job to worker pool
func (p *WorkerPool) Run(fn func()) {
	p.wg.Add(1)
	p.pool <- fn
}

// startPool start worker pool
func (p *WorkerPool) startPool(size int) {
	for i := 0; i < size; i++ {
		go func() {
			for fn := range p.pool {
				fn()
				p.wg.Done()
			}
		}()
	}
}

// Wait wait for all pending jobs are finished
func (p *WorkerPool) Wait() error {
	p.wg.Wait()
	return nil
}

// Close wait for all pending jobs are finished
// and close the channel for good
func (p *WorkerPool) Close() {
	p.wg.Wait()
	close(p.pool)
}
