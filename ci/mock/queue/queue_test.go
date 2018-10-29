package queue

import (
	"testing"
)

const (
	length  = 5
	message = "test"
)

func TestNew(t *testing.T) {
	queue := New(length)
	err := queue.Push([]byte(message))
	if err != nil {
		t.Error("Failed to push message onto queue")
	}
	msgs := queue.Stream()
	select {
	case msg := <-msgs:
		if string(msg) != message {
			t.Errorf("Unexpected message %s, expected %s", string(msg), message)
		}
	default:
		t.Errorf("Could not pop message off the queue.")
	}
}
