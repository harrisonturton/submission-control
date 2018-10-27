package worker

import (
	"github.com/harrisonturton/submission-control/ci/mock/client"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
)

const (
	jobQueueName    = "job_queue"
	resultQueueName = "result_queue"
	queueAddr       = "amqp://guest:guest@localhost:5672/"
	dockerVersion   = "1.38"
)

func TestNew(t *testing.T) {
	client, _ := client.New(dockerVersion)
	jobQueue, _ := queue.New(jobQueueName, queueAddr)
	resultQueue, _ := queue.New(resultQueueName, queueAddr)
	New(jobQueue, resultQueue, client, os.Stdout)
}
