package types

import (
	"encoding/json"
	"net/http"
	"time"
)

// Meta holds metadata about each job.
type Meta struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// AcceptJob is returned when a job is successfully
// marshalled, serialized, and put on the job queue.
type AcceptJob struct {
	Meta
}

// RejectJob is returned when a job cannot be
// marshelled, serialized, or put on the job queue.
type RejectJob struct {
	Message string `json:"message"`
}

// ProcessingJob is returned when a
type ProcessingJob struct {
	Meta
}

// JobResult asda
type JobResult struct {
	Meta
	Output string `json:"output"`
}

// Write will write a response corresponding to the
// AcceptJob instance.
func (a AcceptJob) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(a)
}

// Write will write a response corresponding to the
// ProcessingJob instance.
func (a ProcessingJob) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	return json.NewEncoder(w).Encode(a)
}

// Write will write a response corresponding to the
// RejectJob instance.
func (a RejectJob) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(a)
}

// WriteWith will write a response corresponding to the
// RejectJob instance, with a custom status code.
func (a RejectJob) WriteWith(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(a)
}

// Write will write a response corresponding to the
// JobResult instance.
func (a JobResult) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(a)
}
