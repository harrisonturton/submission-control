package worker

import (
	"bytes"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/queue"
	"github.com/harrisonturton/submission-control/ci/worker/client"
	"io"
	"log"
	"sync"
	"time"
)

type Worker struct {
	Jobs    *queue.Queue
	Results *queue.Queue
	Client  *client.Client
	Logger  *log.Logger
}

// New tries to connect to the job queue, the result queue, and to the docker
// daemon. If any of these fail, an error is returned. Otherwise it returns a
// new Worker instance.
func New(logOut io.Writer, jobQueueName, resultQueueName, addr string) (*Worker, error) {
	jobs, err := queue.New(jobQueueName, addr)
	if err != nil {
		return nil, err
	}
	results, err := queue.New(resultQueueName, addr)
	if err != nil {
		return nil, err
	}
	client, err := client.New("1.38")
	if err != nil {
		return nil, err
	}
	return &Worker{
		Jobs:    jobs,
		Results: results,
		Client:  client,
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

	resp, err := worker.Client.CreateContainer(msg, []string{})
	if err != nil {
		worker.Logger.Printf("Error on create container: %s", err)
		return
	}
	err = worker.Client.WaitForContainer(resp.ID, time.Second*15)
	if err != nil {
		worker.Logger.Printf("Error on wait for container: %s", err)
		return
	}
	logReader, err := worker.Client.ReadContainerLogs(resp.ID, true, true)
	buf := new(bytes.Buffer)
	buf.ReadFrom(logReader)
	worker.Logger.Printf("Container Logs: %s", buf.String())
	worker.Results.Message(fmt.Sprintf("\n%s:\n%s", resp.ID, buf.String()))
}
