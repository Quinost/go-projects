package worker

import (
	"context"
	"fmt"
	"log"
)

type JobProcess[T any] interface {
	Process(data T) error
}

type Job[T any] struct {
	data    T
	process JobProcess[T]
}

type Worker[T any] struct {
	id   int
	jobs <-chan Job[T]
}

func NewJob[T any](data T, process JobProcess[T]) Job[T] {
	return Job[T]{data: data, process: process}
}

func NewWorker[T any](id int, jobs <-chan Job[T]) Worker[T] {
	return Worker[T]{id: id, jobs: jobs}
}

func (w *Worker[T]) Start(ctx context.Context) {
	log.Println("Worker", w.id, "started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker", w.id, "stopping")
			return
		case job, ok := <-w.jobs:
			if !ok {
				log.Println("Worker", w.id, "job channel closed, stopping")
				return
			}

			processName := fmt.Sprintf("%T", job.process)

			log.Println("Worker", w.id, "received job", processName)
			if err := job.process.Process(job.data); err != nil {
				log.Println("Worker", w.id, "Job", processName, "error processing job:", err)
			} else {
				log.Println("Worker", w.id, "completed job", processName)
			}
		}
	}
}
