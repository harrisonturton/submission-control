package store

import (
	"github.com/pkg/errors"
	"time"
)

// WriteUser will upsert data in the database. If the user exists,
// (as defined by the uid), they will get new info. If they do not,
// a new user is created.
func (store *Store) WriteUser(user User) error {
	command := `
INSERT INTO users (uid, first_name, last_name, email, password) VALUES
	($1, $2, $3, $3, $4, $5)
ON CONFLICT DO UPDATE	
`
	_, err := store.db.Exec(command, user.UID, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	return errors.Wrap(err, "failed to write user")
}

// WriteTutorialEnrolment will enrol a user in a tutorial if they are not already enrolled.
func (store *Store) WriteTutorialEnrolment(uid string, tutorialID int) error {
	command := `
INSERT INTO tutorial_enrol (uid, tutorial_id) VALUES
	($1, $2)
ON CONFLICT DO UPDATE	
`
	_, err := store.db.Exec(command, uid, tutorialID)
	return errors.Wrap(err, "failed to write tutorial enrolment")
}

// WriteCourseEnrolment will enrol a user in a course if they are not already enrolled
func (store *Store) WriteCourseEnrolment(uid string, courseID int, role int) error {
	command := `
INSERT INTO enrol (role, uid, course_id) VALUES
	($1, $2, $3)
ON CONFLICT DO UPDATE`
	_, err := store.db.Exec(command, role, uid, courseID)
	return errors.Wrap(err, "failed to write course enrolment")
}

// WriteSubmission will write a submission to the database.
func (store *Store) WriteSubmission(uid string, assessmentID int, title string, description string, file []byte) error {
	command := `
INSERT INTO submissions (uid, timestamp, assessment_id, data, title, description, feedback) VALUES
	($1, $2, $3, $4, $5, $6, '')
`
	_, err := store.db.Exec(command, uid, time.Now(), assessmentID, file, title, description)
	return errors.Wrap(err, "failed to write submission")
}

// WriteSubmissionFeedback will write new feedback onto a submission
func (store *Store) WriteSubmissionFeedback(submissionID int, feedback string) error {
	command := `UPDATE submissions SET feedback=$1 WHERE id=$2`
	_, err := store.db.Exec(command, feedback, submissionID)
	return errors.Wrap(err, "failed to write submission feedback")
}

// WriteTestResult will write a test result to the database.
func (store *Store) WriteTestResult(submissionID int, testWarnings, testErrors, resultType string) error {
	command := `
INSERT INTO test_results (submission_id, warnings, errors, type) VALUES
	($1, $2, $3, $4)
`
	var resultTypeID int
	switch resultType {
	case "success":
		resultTypeID = 1
		break
	case "warnings":
		resultTypeID = 2
		break
	case "failed":
		resultTypeID = 3
		break
	case "untested":
		resultTypeID = 4
		break
	}
	_, err := store.db.Exec(command, submissionID, testWarnings, testErrors, resultTypeID)
	return errors.Wrap(err, "failed to write test result")
}
