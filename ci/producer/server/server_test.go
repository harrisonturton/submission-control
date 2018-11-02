// +build unit

package server

import (
	"github.com/harrisonturton/submission-control/ci/cache"
	"github.com/harrisonturton/submission-control/ci/mock/queue"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	goodAll       = `{"version":"1","environment":{"image":"hello-world","vars":{"key":"value"}}}`
	goodNoVersion = `{"environment":{"image":"haskell"}}`
	badNoEnv      = `{"version":"1"}`
	badNoImage    = `{"version":"1","environment":{"vars":{"k":"v"}}}`
)

const addr = "localhost:8080"

func TestHandle(t *testing.T) {
	jobs := queue.New(5)
	cache := cache.New(5, time.Hour*10)
	server := New(os.Stdout, jobs, cache, addr)
	checkResponse := func(body string, statusCode int) {
		reader := strings.NewReader(body)
		req := httptest.NewRequest("POST", addr, reader)
		resp := httptest.NewRecorder()
		server.handleRequest(resp, req)
		if resp.Code != statusCode {
			t.Errorf("Expected status code %d, but got %d", statusCode, resp.Code)
		}
		if resp.Header()["Content-Type"][0] != "application/json" {
			t.Errorf("Expected Content-Type application/json, but got %s", resp.Header()["Content-Type"][0])
		}
	}
	checkResponse(goodAll, 201)
	checkResponse(goodNoVersion, 201)
	checkResponse(badNoEnv, 400)
	checkResponse(badNoImage, 400)
}
