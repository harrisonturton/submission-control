package routes

import (
	"bytes"
	"github.com/harrisonturton/submission-control/store/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	// Make the request & record the response
	resp := httptest.NewRecorder()
	body, err := json.Marshal(LoginRequest{
		UID:      "u6386433",
		Password: "testing123",
	})
	req := http.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
	authHandler.ServeHTTP(resp, req)
}
