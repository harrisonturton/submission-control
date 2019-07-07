package store

import (
	"database/sql"
	"github.com/pkg/errors"
)

// Reader mimics the Store query methods, allowing
// us to mock the database.
type Reader interface {
	GetUser(uid string) (*User, error)
	GetAssessment(uid string) ([]Assessment, error)
	GetSubmissions(uid string) ([]Submission, error)
	GetEnrolment(uid string) ([]Enrolment, error)
	GetTutorialEnrolment(uid string) ([]TutorialEnrolment, error)

	GetCourse(courseID int) (*Course, error)
	GetTutorial(tutorialID int) (*Tutorial, error)
}

// Writer mimics the Store write methods, allowing
// us to mock the database.
type Writer interface {
	WriteUser(user User) error
	WriteTutorialEnrolment(uid string, tutorialID int) error
	WriteCourseEnrolment(uid string, courseID int, role int) error
	WriteSubmission(uid string, assessmentID int, title string, description string, file []byte) (int64, error)
	WriteTestResult(testWarnings, testErrors, resultType string) (int64, error)
}

// Store represents the database. It does NOT
// represent a single connection, since database/sql
// manages a connection pool beneath the hood.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance, and ensures
// the database can be connected to.
func NewStore(db *sql.DB) (*Store, error) {
	err := db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return &Store{db}, nil
}
