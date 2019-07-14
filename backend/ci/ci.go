package ci

import (
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"os"
	"sync"
)

// Ci is the automated testing service
type Ci struct {
	logger *log.Logger
	store  *store.Store
	jobs   chan Job
}

// Job is the job each worker works with
type Job struct {
	SubmissionID int
	AssessmentID int
}

// NewCi creates a new testing service instance
func NewCi(logger *log.Logger, store *store.Store) *Ci {
	return &Ci{
		logger: logger,
		store:  store,
		jobs:   make(chan Job, 500),
	}
}

// Run will spawn a number of works all listening
// for work on the job channel
func (ci *Ci) Run(numWorkers int, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	err := os.MkdirAll(testDir, 0777)
	if err != nil {
		ci.logger.Printf("Failed to create testing dir: %v\n", err)
		return
	}
	ci.logger.Printf("Starting %d workers...\n", numWorkers)
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(ci.logger, ci.store)
		go worker.Run(ci.jobs)
	}
	<-done
	close(ci.jobs)
}

// Notify that a new submission has been added
func (ci *Ci) Notify() {
	submissions, err := ci.store.GetUntestedSubmissionIDs()
	if err != nil {
		ci.logger.Printf("Failed to get untested submission IDs in CI: %v\n", err)
		return
	}
	for _, submission := range submissions {
		ci.jobs <- Job{submission.ID, submission.AssessmentID}
	}
}
