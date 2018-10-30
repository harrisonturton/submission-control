// +build unit
package client

import (
	"testing"
	"time"
)

const (
	imageID = "test_image_id"
)

func TestCreate(t *testing.T) {

	client := New()
	id, err := client.Create(imageID)
	if err != nil {
		t.Error("Could not create container")
	}
	if id == "" {
		t.Error("Blank container ID on create")
	}

	err = client.Start(id)
	if err != nil {
		t.Error("Error starting container")
	}

	err = client.Start(id)
	if err == nil {
		t.Error("Expected error when trying to start running container, but recieved none!")
	}

	logs, err := client.ReadLogs(id, true, true)
	if err != nil {
		t.Error("Error reading container logs")
	}
	if logs == "" {
		t.Error("Container logs are blank")
	}

	err = client.Wait(id, time.Second*5)
	if err != nil {
		t.Error("Container wait returned an error")
	}

	if !client.Healthy() {
		t.Error("Unhealthy client")
	}

}
