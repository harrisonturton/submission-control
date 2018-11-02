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

func TestStatusCodes(t *testing.T) {
	expectStatus := func(body string, statusCode int) {
		reader := strings.NewReader(body)
		resp, err := http.Post(addr, encoding, reader)
		if err != nil {
			t.Errorf("Failed to post: %s", err)
		}
		if resp.StatusCode != statusCode {
			t.Errorf("Expected status code %d, but got %d", statusCode, resp.StatusCode)
		}
	}
	expectStatus(goodAll, 201)
	expectStatus(goodNoVersion, 201)
	expectStatus(badNoEnv, 400)
	expectStatus(badNoImage, 400)
}
