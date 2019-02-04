package routes

import "github.com/harrisonturton/submission-control/backend/store"

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
	UID      string `json:"uid"`
	Password string `json:"password"`
}

// RefreshRequest is a POST request sent to the /refresh
// endpoint. If the client has a valid JWT token, it allows
// them to refresh it, allowing for persistent logins.
type RefreshRequest struct {
	Token string `json:"token"`
}

// TokenResponse responds with a JWT token.
type TokenResponse struct {
	StatusCode int    `json:"status"`
	Token      string `json:"token"`
}

// EnrolResponse is the response to an EnrolledRequest
type EnrolResponse struct {
	StatusCode int            `json:"status"`
	Courses    []store.Course `json:"courses"`
}

// AssessmentResponse is the response to an AssessmentRequest
type AssessmentResponse struct {
	StatusCode int                `json:"status"`
	Assessment []store.Assessment `json:"assessment"`
}

// UserResponse is sent in response to a GET request on /user
type UserResponse struct {
	StatusCode int        `json:"status"`
	User       store.User `json:"user"`
}

// SubmissionResponse is the response sent when replying to a
// submission request on the /submission endpoint
type SubmissionResponse struct {
	StatusCode  int                `json:"status"`
	Submissions []store.Submission `json:"submissions"`
}

// StudentStateResponse contains all the state data required to render
// a students page on the client.
type StudentStateResponse struct {
	User        store.User         `json:"user"`
	Assessment  []store.Assessment `json:"assessment"`
	Submissions []store.Submission `json:"submissions"`
	Enrolled    []store.Enrolment  `json:"enrolled"`
}
