// +build unit

package server

import (
	"github.com/harrisonturton/submission-control/ci/cache"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"os"
	"testing"
	"time"
)

const (
	serverAddr = "localhost:8080"
)

func TestServe(t *testing.T) {
	jobs := queue.New(5)
	cache := cache.New(5, time.Hour*10)
	New(os.Stdout, jobs, cache, serverAddr)
}
