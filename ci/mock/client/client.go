package client

import (
	"fmt"
	"math/rand"
	"time"
)

// Client is a mocked instance of a connection
// to the Docker daemon.
type Client struct {
	Containers []string
	Running    map[string]bool
}

// New creates a new instance of the mocked client.
func New() *Client {
	return &Client{
		Containers: []string{},
		Running:    map[string]bool{},
	}
}

// Healthy will always return true, because it is mocked.
func (client *Client) Healthy() bool {
	return true
}

// Create creates a new mock container, and returns an random ID.
func (client *Client) Create(imageID string) (string, error) {
	id := randString(10)
	client.Containers = append(client.Containers, id)
	return id, nil
}

// Start launches a non-running container.
func (client *Client) Start(containerID string) error {
	if client.Running[containerID] {
		return fmt.Errorf("container %s already running", containerID)
	}
	return nil
}

// ReadLogs reads some fake logs from a given container.
func (client *Client) ReadLogs(containerID string, showStdout bool, showStderr bool) (string, error) {
	return fmt.Sprintf("Logs for %s.\nStdout: %t | Stderr: %t", containerID, showStdout, showStderr), nil
}

// Wait doesn't actually wait, but will just return instantly.
func (client *Client) Wait(containerID string, timeout time.Duration) error {
	return nil
}

// randString creates a string of random letters of a certain length
func randString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
