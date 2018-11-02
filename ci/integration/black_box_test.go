// +build integration

package integration

import (
	"net/http"
	"strings"
	"testing"
)

const (
	goodAll       = `{"version":"1","environment":{"image":"hello-world","vars":{"key":"value"}}}`
	goodNoVersion = `{"environment":{"image":"haskell"}}`
	badNoEnv      = `{"version":"1"}`
	badNoImage    = `{"version":"1","environment":{"vars":{"k":"v"}}}`
)

const addr = "http://localhost:8080/"
const encoding = "application/json"

func Test(t *testing.T) {
	client := &http.Client{}
	reader := strings.NewReader(badNoEnv)
	req, err := http.NewRequest("POST", addr, reader)
	if err != nil {
		t.Fatalf("Failed to created request: %s", err)
		return
	}
	req.Close = true
	req.Header.Set("Content-Type", encoding)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
		return
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected code %d, but got %d", 400, resp.StatusCode)
	}
}

/*func Test(t *testing.T) {
	client := &http.Client{}
	expectStatus := func(body string, statusCode int) {
		req, err := newRequest(goodAll)
		if err != nil {
			t.Errorf("Failed to create request for %s: %s", body, err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Failed to get response for %s: %s", body, err)
			return
		}
		if resp.StatusCode != statusCode {
			t.Errorf("Expected status code %d, got %d for body %s", statusCode, resp.StatusCode, body)
		}
	}
	expectStatus(goodAll, 201)
	expectStatus(goodNoVersion, 201)
	expectStatus(badNoEnv, 400)
	expectStatus(badNoImage, 400)
}*/

func newRequest(body string) (*http.Request, error) {
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", addr, reader)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Content-Type", encoding)
	return req, nil
}
