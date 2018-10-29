package client

import "testing"

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

	logs, err := client.ReadLogs(id, true, true)
	if err != nil {
		t.Error("Error reading container logs")
	}
	if logs == "" {
		t.Error("Container logs are blank")
	}

	if !client.Healthy() {
		t.Error("Unhealthy client")
	}

}
