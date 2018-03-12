package usecases

import (
	"fmt"

	// "github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type Worker struct {
	jobQueue chan Job
	done     chan bool
	logger   domain.Logger
}

func (w *Worker) Start() {
	for {
		select {
		case job := <-w.jobQueue:
			err := job.DoWork()
			if err != nil {
				w.logger.LogError(fmt.Sprintf("Worker returned error for job: Error:: %v", err))
			}
		case <-w.done:
			return
		}
	}
}

func NewWorker(jobQueue chan Job, closeChannel chan bool) *Worker {
	w := Worker{
		jobQueue: jobQueue,
		done:     closeChannel,
	}

	return &w
}

type Dispatcher struct {
	Worker       *Worker
	JobQueue     chan Job
	CloseChannel chan bool
}

func (d *Dispatcher) Run(maxWorkers int) {
	for i := 1; i <= maxWorkers; i++ {
		go d.Worker.Start()
	}
}

func (d *Dispatcher) Stop() {
	go func() {
		d.CloseChannel <- true
	}()
}

func NewDispatcher(jobQueue chan Job, closeChannel chan bool) *Dispatcher {
	d := Dispatcher{
		Worker:       NewWorker(jobQueue, closeChannel),
		JobQueue:     jobQueue,
		CloseChannel: closeChannel,
	}

	return &d
}
