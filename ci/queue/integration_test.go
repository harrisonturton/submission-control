// +build integration

package queue

import (
	"testing"
	"time"
)

const (
	queueName = "test_queue"
	host      = "localhost"
	queueAddr = "amqp://guest:guest@" + host + ":5672/"
	message   = "test message"
	timeout   = 5 * time.Second
)

func Test(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// Setup queue
	queue, err := New(queueName, queueAddr)
	if err != nil {
		t.Fatalf("Failed with error: %s", err)
	}

	t.Run("Push and pop", func(t *testing.T) {
		if err := queue.Push([]byte(message)); err != nil {
			t.Fatalf("Failed with error: %s", err)
		}
		msgs := queue.Stream()
		var response []byte
		select {
		case msg := <-msgs:
			response = msg
		case <-time.After(timeout):
			t.Fatalf("Failed to pop message: timed out")
		}
		if string(response) != message {
			t.Fatalf("Recieved unexpected message: %s", string(response))
		}
	})
	err = queue.Close()
	if err != nil {
		t.Error("Failed to close the connection")
	}
}
