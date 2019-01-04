package routes

// LoginRequest is a POST request sent to the /auth
// endpoint. The server responds with a time-limited
// JWT access token, which must be attached to the
// Authentication header for all consequent API requests.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the data sent in response to a login
// request.
type LoginResponse struct {
	Token string `json:"token"`
}

// RefreshRequest is a POST request sent to the /refresh
// endpoint. If the client has a valid JWT token, it allows
// them to refresh it, allowing for persistent logins.
type RefreshRequest struct {
	Token string `json:"token"`
}
