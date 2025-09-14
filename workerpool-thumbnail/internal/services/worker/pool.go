package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool[T any] struct {
	workerCount int
	jobQueue    chan Job[T]
	workers     []*Worker[T]
	wg          sync.WaitGroup
}

func NewPool[T any](workerCount int, jobQueueSize int) *Pool[T] {
	pool := &Pool[T]{
		workerCount: workerCount,
		jobQueue:    make(chan Job[T], jobQueueSize),
		workers:     make([]*Worker[T], 0, workerCount),
	}

	for i := range workerCount {
		worker := NewWorker(i+1, pool.jobQueue)
		pool.workers = append(pool.workers, &worker)
	}

	return pool
}

func (p *Pool[T]) Start(ctx context.Context) {
	log.Printf("Starting worker pool with %d workers", p.workerCount)
	for _, worker := range p.workers {
		p.wg.Add(1)
		go func(w *Worker[T]) {
			defer p.wg.Done()
			w.Start(ctx)
		}(worker)
	}
}

func (p *Pool[T]) Stop() {
	processName := fmt.Sprintf("%T", p)
	log.Println("Stopping worker pool", processName)
	close(p.jobQueue)
	p.wg.Wait()
}

func (p *Pool[T]) Submit(job Job[T]) {
	p.jobQueue <- job
}
