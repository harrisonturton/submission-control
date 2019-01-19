package store

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

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

// GetUserByEmail will return a single user with a
// matching email address.
func (store *Store) GetUserByEmail(email string) (*User, error) {
	var firstname, lastname, uid string
	var passwordHash []byte
	query := "SELECT first_name, last_name, password, uid FROM users WHERE email = $1"
	err := store.db.QueryRow(query, email).Scan(&firstname, &lastname, &passwordHash, &uid)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch account")
	}
	return &User{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		PasswordHash: string(passwordHash),
		UID:          uid,
	}, nil
}

// GetCoursesByUser will return an array of all the courses a user is
// enrolled in.
func (store *Store) GetCoursesByUser(uid string) ([]Course, error) {
	// Get a list of CourseIDs
	query := "SELECT course_id FROM enrol WHERE user_uid = $1"
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Course{}, err
	}
	var courseIDs []string
	for rows.Next() {
		var courseID string
		err := rows.Scan(&courseID)
		if err != nil {
			log.Println(err.Error())
			log.Println("Failed in rows.Next for " + courseID)
			continue
		}
		courseIDs = append(courseIDs, courseID)
	}
	rows.Close()
	// Use the IDs to query for the actual course items
	var courses []Course
	for _, courseID := range courseIDs {
		var name, courseCode, period string
		var year int
		var id int
		query = "SELECT id, name, course_code, period, year FROM courses WHERE id = $1"
		err := store.db.QueryRow(query, courseID).Scan(&id, &name, &courseCode, &period, &year)
		if err != nil {
			log.Println(err.Error())
			log.Printf("Failed in queryRow for %d\n", courseID)
			continue
		}
		courses = append(courses, Course{
			ID:         id,
			Name:       name,
			CourseCode: courseCode,
			Period:     period,
			Year:       year,
		})
	}
	return courses, nil
}

// GetAssessmentForUser will fetch a list of all assessments (for all courses) for a single user.
func (store *Store) GetAssessmentForUser(uid string) ([]Assessment, error) {
	query := `
SELECT
	assessment.id as assessment_id,
	assessment.course_id,
	name,
	type
FROM
	assessment JOIN enrol ON assessment.course_id = enrol.course_id
WHERE enrol.user_uid = $1;
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Assessment{}, nil
	}
	var assessment []Assessment
	for rows.Next() {
		var assessmentID, courseID int
		var name, assType string
		err := rows.Scan(&assessmentID, &courseID, &name, &assType)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		assessment = append(assessment, Assessment{
			ID:       assessmentID,
			Name:     name,
			Type:     assType,
			CourseID: courseID,
		})
	}
	rows.Close()
	return assessment, nil
}

// GetSubmissionsForUser will return all the submissions made
// by the user to every assessment they've had.
func (store *Store) GetSubmissionsForUser(uid string) ([]Submission, error) {
	//query := "SELECT id, assessment_id, uid, title, description, feedback FROM submissions WHERE uid = $1"
	query := `
SELECT
	assessment.course_id, 
	assessment.id  AS assessment_id,
	submissions.id AS submission_id, 
	submissions.title,
	submissions.description,
	feedback
FROM assessment
JOIN submissions ON assessment.id = submissions.assessment_id
WHERE uid = $1;
	`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Submission{}, nil
	}
	var submissions []Submission
	for rows.Next() {
		var courseID, assessmentID, submissionID int
		var title, description, feedback string
		err := rows.Scan(&courseID, &assessmentID, &submissionID, &title, &description, &feedback)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		submissions = append(submissions, Submission{
			ID:           submissionID,
			AssessmentID: assessmentID,
			CourseID:     courseID,
			Title:        title,
			Description:  description,
			Feedback:     feedback,
		})
	}
	rows.Close()
	return submissions, nil
}
