package worker

import (
	"fmt"
	"github.com/harrisonturton/submission-control/ci/client"
	"github.com/harrisonturton/submission-control/ci/queue"
	"io"
	"log"
	"sync"
	"time"
)

// Worker represents the processes that take student projects from
// the test queue, run them in a container, and puts the results
// on the result queue.
type Worker struct {
	Jobs    queue.Queue
	Results queue.Queue
	Client  client.Client
	Logger  *log.Logger
}

// New creates a new Worker.
func New(jobs queue.Queue, results queue.Queue, client client.Client, logOut io.Writer) *Worker {
	return &Worker{
		Jobs:    jobs,
		Results: results,
		Client:  client,
		Logger:  log.New(logOut, "", log.LstdFlags),
	}
}

// Run will continuously pop tasks off the queue, run them inside
// a container, and push the results onto the result queue.
func (worker *Worker) Run(done chan bool, wg *sync.WaitGroup) {
	worker.Logger.Printf("Waiting for jobs on queue %s", worker.Jobs.Name())
	worker.Jobs.Consume(wg, done, worker.handleJob)
}

// handleJob is called for every task that is recieved from the
// job queue. It runs the code inside a container, and puts the
// STDOUT on the results queue.
func (worker *Worker) handleJob(msg string) {
	worker.Logger.Printf("Recieved job: %s", msg)
	id, err := worker.Client.Create(msg)
	if err != nil {
		worker.Logger.Printf("Error on create container: %s", err)
		return
	}
	err = worker.Client.Wait(id, time.Second*5)
	if err != nil {
		worker.Logger.Printf("Error on wait for container: %s", err)
		return
	}
	logs, err := worker.Client.ReadLogs(id, true, true)
	if err != nil {
		worker.Logger.Printf("Error on fetching container logs: %s", err)
		return
	}
	worker.Logger.Printf("Container Logs: %s", logs)
	worker.Results.Message(fmt.Sprintf("\n%s:\n%s", id, logs))
}
