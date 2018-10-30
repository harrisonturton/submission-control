// +build unit
package server

import (
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
)

const (
	serverAddr = "localhost:8080"
)

func TestServe(t *testing.T) {
	jobs := queue.New(5)
	New(os.Stdout, jobs, serverAddr)
}
