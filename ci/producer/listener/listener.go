package listener

import (
	"github.com/harrisonturton/submission-control/test-engine/queue"
	"io"
	"log"
	"sync"
)

type Listener struct {
	Results *queue.Queue
	Logger  *log.Logger
}

// New tries to connect to the result queue. If successful, it will return
// a Listener instance. Otherwise, it will fail with an error.
func New(logOut io.Writer, resultQueueName, addr string) (*Listener, error) {
	results, err := queue.New(resultQueueName, addr)
	if err != nil {
		return nil, err
	}
	return &Listener{
		Results: results,
		Logger:  log.New(logOut, "", log.LstdFlags),
	}, nil
}

// Run will continuouslly pop messages off the result queue,
// and handle them with listener.handleResult
func (listener *Listener) Run(done chan bool, wg *sync.WaitGroup) {
	listener.Logger.Printf("Waiting for results on queue %s", listener.Results.Queue.Name)
	listener.Results.Consume(wg, done, listener.handleResult)
}

// handleResult processes every message that the server recieves
// from the result queue.
func (listener *Listener) handleResult(message string) {
	listener.Logger.Printf("Recieved result: %s", message)
}
