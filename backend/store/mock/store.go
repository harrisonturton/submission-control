package mock

import (
	"github.com/harrisonturton/submission-control/backend/store"
)

// Store mocks the real store implementation in the
// parent directory. Rather than connecting to the
// database, it returns dummy data.
type Store struct{}

// NewStore creates a new mock store.
func NewStore() (*Store, error) {
	return &Store{}, nil
}

// GetUser will find a user for the given uid.
func (store *Store) GetUser(uid string) (store.User, error) {
	return store.User{
		UID:          uid,
		Email:        uid + "@anu.edu.au",
		FirstName:    "Jonathan",
		LastName:     "Smith",
		PasswordHash: "dummy-password-hash",
	}, nil
}

// GetAssessment will find a list of assessments for the given user.
func (store *Store) GetAssessment(uid string) ([]store.Assessment, error) {
	return []store.Assessment{
		// Course 1
		Assessment{
			ID:       1,
			CourseID: 1,
			Name:     "Gliding In Space",
			Type:     "assignment",
		},
		Assessment{
			ID:       2,
			CourseID: 1,
			Name:     "Distributed Server",
			Type:     "lab",
		},
		// Course 2
		Assessment{
			ID:       3,
			CourseID: 2,
			Name:     "Musical Instrument",
			Type:     "assignment",
		},
		Assessment{
			ID:       4,
			CourseID: 2,
			Name:     "Blinking LED",
			Type:     "lab",
		},
	}, nil
}

// GetSubmissions will find a list of submissions for the given user.
func (store *Store) GetSubmissions(uid string) ([]store.Submission, error) {
	return []store.Submission{
		Submission{
			ID:           1,
			AssessmentID: 1,
			CourseID:     1,
			Title:        "First attempt at Gliding In Space",
			Description:  "Hopefully this works. I think it'll pass most tests",
			Feedback:     "Nice work - just check those last couple of tests! Comments look good, but don't comment everything.",
		},
	}, nil
}

// GetEnrolment will find a list of enrolments for the given user.
func (store *Store) GetEnrolment(uid string) ([]store.Enrolment, error) {
	return []store.Enrolment{
		Enrolment{
			Course: Course{
				ID:         1,
				Name:       "Concurrent & Distributed Programming",
				CourseCode: "COMP2310",
				Period:     store.PeriodFirst,
				Year:       2018,
			},
			UID:  uid,
			Role: "student",
		},
		Enrolment{
			Course: Course{
				ID:         2,
				Name:       "Computer Organisation & Execution",
				CourseCode: "COMP2300",
				Period:     store.PeriodSecond,
				Year:       2017,
			},
			UID:  uid,
			Role: "tutor",
		},
	}, nil
}
