package types

import (
	"net/http"
	"time"
)

// Meta holds metadata about each job.
type Meta struct {
	ID        string
	Timestamp time.Time
}

// AcceptJob is returned when a job is successfully
// marshalled, serialized, and put on the job queue.
type AcceptJob struct {
	Meta
}

// RejectJob is returned when a job cannot be
// marshelled, serialized, or put on the job queue.
type RejectJob struct {
	Message string
}

// ProcessingJob is returned when a
type ProcessingJob struct {
	Meta
}

// JobResult asda
type JobResult struct {
	Meta
	Output string
}

// ResponseCreated will write a Status 201 response
func ResponseCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status_code":201,"status":"201 CREATED"}`))
}

// ResponseAccepted will write a Status 202 response
func ResponseAccepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status_code":202,"status":"202 ACCEPTED"}`))
}

// ResponseBadRequest will write a Status 400 response
func ResponseBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"status_code":400,"status":"400 BAD REQUEST"}`))
}

// ResponseInternalServerError will write a Status 500 response
func ResponseInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"status_code":500,"status":"500 INTERNAL SERVER ERROR"}`))
}
