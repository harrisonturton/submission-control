package worker

import (
	"github.com/harrisonturton/submission-control/test-engine/queue"
	"io"
	"log"
	"sync"
)

type Worker struct {
	Jobs    *queue.Queue
	Results *queue.Queue
	Logger  *log.Logger
}

// New tries to connect to both the job and result queues. If successful,
// it will return a Worker instance. Otherwise, it will fail with an error.
func New(logOut io.Writer, jobQueueName, resultQueueName, addr string) (*Worker, error) {
	jobs, err := queue.New(jobQueueName, addr)
	if err != nil {
		return nil, err
	}
	results, err := queue.New(resultQueueName, addr)
	if err != nil {
		return nil, err
	}
	return &Worker{
		Jobs:    jobs,
		Results: results,
		Logger:  log.New(logOut, "", log.LstdFlags),
	}, nil
}

// Run will continuously pop messages off the queue, process the
// container, and put the results into the results queue.
func (worker *Worker) Run(done chan bool, wg *sync.WaitGroup) {
	worker.Logger.Printf("Waiting for jobs on queue %s", worker.Jobs.Queue.Name)
	worker.Jobs.Consume(wg, done, worker.handleJob)
}

// handleJob is called every time a job is recieved from the job
// queue.
func (worker *Worker) handleJob(msg string) {
	worker.Logger.Printf("Recieved job: %s", msg)
	worker.Results.Message("Handled " + msg)
}
