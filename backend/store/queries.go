package store

import (
	"github.com/pkg/errors"
	"log"
)

// GetUser will return a single user with a matching uid.
func (store *Store) GetUser(uid string) (*User, error) {
	var firstname, lastname, email string
	var passwordHash []byte
	query := "SELECT first_name, last_name, password, email FROM users WHERE uid = $1"
	err := store.db.QueryRow(query, uid).Scan(&firstname, &lastname, &passwordHash, &email)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch account")
	}
	enrolment, err := store.GetEnrolment(uid)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch enrolment for account")
	}
	return &User{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		PasswordHash: string(passwordHash),
		UID:          uid,
		Enrolment:    enrolment,
	}, nil
}

// GetUserRole will fetch the role for a user who is enrolled in a course.
// It will return an error if the user is not enrolled in the course, or
// if the user cannot be found.
func (store *Store) GetUserRole(uid, courseID string) (string, error) {
	query := "SELECT role FROM enrol WHERE user_uid = $1 AND course_id $2"
	var role string
	err := store.db.QueryRow(query, uid, courseID).Scan(&role)
	if err != nil {
		return "", errors.Wrap(err, "could not fetch role")
	}
	return role, nil
}

// GetAssessment will fetch a list of all assessments (for all courses) for a single user.
func (store *Store) GetAssessment(uid string) ([]Assessment, error) {
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

// GetSubmissions will return all the submissions made
// by the user to every assessment they've had.
func (store *Store) GetSubmissions(uid string) ([]Submission, error) {
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

// GetEnrolment will return all the enrolments (mappings from user
// to course, with a role) for a single user.
func (store *Store) GetEnrolment(uid string) ([]Enrolment, error) {
	query := `
SELECT id, user_uid as uid, role, name, course_code, period, year
FROM enrol
JOIN courses ON enrol.course_id = courses.id
WHERE enrol.user_uid = $1;
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Enrolment{}, err
	}
	var enrolment []Enrolment
	for rows.Next() {
		var id, year int
		var uid, name, role, courseCode, period string
		err := rows.Scan(&id, &uid, &role, &name, &courseCode, &period, &year)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		enrolment = append(enrolment, Enrolment{
			Course: Course{
				ID:         id,
				Name:       name,
				CourseCode: courseCode,
				Period:     period,
				Year:       year,
			},
			UID:  uid,
			Role: role,
		})
	}
	rows.Close()
	return enrolment, nil
}
