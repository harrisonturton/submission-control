package ci

import (
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	testDir            = "/tmp/submission-app-test/"
	testDirPrefix      = "testing-submission-"
	specFilename       = "project.spec.sh"
	submissionFilename = "project.zip"
)

// Worker takes jobs off the job channel
// and does the automated testing.
type Worker struct {
	log   *log.Logger
	store *store.Store
}

// NewWorker creates a new worker instance
func NewWorker(log *log.Logger, store *store.Store) *Worker {
	return &Worker{
		log:   log,
		store: store,
	}
}

// Run will handle a job when it comes onto the channel
func (worker *Worker) Run(jobs <-chan Job) {
	for job := range jobs {
		worker.handleJob(job.SubmissionID, job.AssessmentID)
	}
}

func (worker *Worker) handleJob(submissionID int, assessmentID int) {
	worker.log.Printf("handling { submission %d, assessment %d }\n", submissionID, assessmentID)
	err := worker.writeFiles(assessmentID, submissionID)
	if err != nil {
		worker.log.Printf("failed to write files: %v\n", err)
		return
	}
	stdout, err := worker.runTestSpec(submissionID)
	if err != nil {
		worker.log.Printf("failed to run test spec: %v\n", err)
		return
	}
	worker.log.Printf("Tested: %s\n", string(stdout))
	//err := worker.store.WriteTestResult(submissionID, "", "", "success")
	//if err != nil {
	//	worker.log.Printf("failed to save test result: %v\n", err)
	//}
}

func (worker *Worker) runTestSpec(submissionID int) ([]byte, error) {
	path := writeTestDir(submissionID, specFilename)
	out, err := exec.Command("sh", path).Output()
	return out, err
}

func (worker *Worker) writeFiles(assessmentID, submissionID int) error {
	err := os.MkdirAll(writeTestDir(submissionID, ""), 0777)
	if err != nil {
		worker.log.Printf("Failed to create testing dir for submission: %v\n", err)
		return err
	}
	err = worker.writeSubmissionFiles(submissionID)
	if err != nil {
		return err
	}
	err = worker.writeTestSpec(assessmentID, submissionID)
	if err != nil {
		return err
	}
	return nil
}

// writeSubmissionFiles will save the submission .zip file to the testing folder
func (worker *Worker) writeSubmissionFiles(submissionID int) error {
	data, err := worker.store.GetSubmissionFiles(submissionID)
	if err != nil {
		worker.log.Printf("Failed to get submission %d files in worker: %v\n", submissionID, err)
		return err
	}
	file, err := os.Create(writeTestDir(submissionID, submissionFilename))
	defer file.Close()
	if err != nil {
		worker.log.Printf("Failed to write submission files to temp dir for testing: %d: %v\n", submissionID, err)
		return err
	}
	file.Write(data)
	return nil
}

// writeTestSpec will save the assessment test spec bash file to the testing folder
func (worker *Worker) writeTestSpec(assessmentID int, submissionID int) error {
	spec, err := worker.store.GetAssessmentSpec(assessmentID)
	if err != nil {
		worker.log.Printf("Failed to get assessment spec %d in worker: %v\n", assessmentID, err)
		return err
	}
	specFile, err := os.Create(writeTestDir(submissionID, specFilename))
	defer specFile.Close()
	if err != nil {
		worker.log.Printf("Failed to write submission spec to temp dir for testing: %d: %v\n", submissionID, err)
		return err
	}
	err = specFile.Chmod(0700)
	if err != nil {
		worker.log.Printf("Failed to give 0700 permissions to spec file for submission %d: %v\n", submissionID, err)
		return err
	}
	specFile.Write(spec)
	return nil
}

func writeTestDir(submissionID int, file string) string {
	return testDir + testDirPrefix + strconv.Itoa(submissionID) + "/" + file
}
