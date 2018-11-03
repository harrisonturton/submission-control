// +build unit

package worker

import (
	"github.com/harrisonturton/submission-control/ci/mock/client"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"github.com/harrisonturton/submission-control/ci/types"
	"os"
	"testing"
	"time"
)

func TestHandleJob(t *testing.T) {
	// Create worker
	client := client.New()
	jobs := queue.New(5)
	results := queue.New(5)
	worker := New(jobs, results, client, os.Stdout)
	// Create job
	testJob := types.TestJob{
		Timestamp: time.Now(),
		Config: types.TestConfig{
			Version: "1",
			Env: types.TestEnv{
				Image: "hello-world",
			},
		},
	}
	// Test job
	worker.handleJob(testJob)
	// Ensure something is placed on result queue
	if results.Length() != 1 {
		t.Errorf("Expected queue length to be %d, but got %d", 1, results.Length())
	}
}
