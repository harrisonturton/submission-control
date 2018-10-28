package server

import (
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
)

const (
	jobQueue   = "result_queue"
	queueAddr  = "amqp://guest:guest@localhost:5672/"
	serverAddr = "localhost:8080"
)

func TestServe(t *testing.T) {
	jobs, _ := queue.New(jobQueue, queueAddr)
	New(os.Stdout, jobs, serverAddr)
}
