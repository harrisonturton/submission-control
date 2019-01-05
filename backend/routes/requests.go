package routes

// FailedResponse is a common response for failed API requests.
type FailedResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

// LoginRequest is a POST request sent to the /auth
// endpoint. The server responds with a time-limited
// JWT access token, which must be attached to the
// Authentication header for all consequent API requests.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest is a POST request sent to the /refresh
// endpoint. If the client has a valid JWT token, it allows
// them to refresh it, allowing for persistent logins.
type RefreshRequest struct {
	StatusCode int    `json:"status"`
	Token      string `json:"token"`
}

// TokenResponse responds with a JWT token.
type TokenResponse struct {
	StatusCode int    `json:"status"`
	Token      string `json:"token"`
}
