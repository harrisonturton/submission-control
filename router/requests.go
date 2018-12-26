package router

// LoginRequest is a POST request sent to the /auth
// endpoint. The server responds with a time-limited
// JWT access token, which must be attached to the
// Authentication header for all consequent API requests.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
