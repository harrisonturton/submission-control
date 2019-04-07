package store

import (
	"github.com/pkg/errors"
	"log"
	"time"
)

// GetUser will fetch the data for a single user
func (store *Store) GetUser(uid string) (*User, error) {
	query := `
SELECT
	first_name,
	last_name,
	email,
	password
FROM users
WHERE uid = $1
`
	var firstname, lastname, email string
	var passwordHash []byte
	err := store.db.
		QueryRow(query, uid).
		Scan(&firstname, &lastname, &email, &passwordHash)
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
		UID:          uid,
		Enrolment:    enrolment,
		PasswordHash: string(passwordHash),
	}, nil
}

// GetEnrolment will fetch all the enrolled courses for a given user
func (store *Store) GetEnrolment(uid string) ([]Enrolment, error) {
	query := `
SELECT
	courses.id,
	roles.role,
	name,
	code,
	periods.period,
	year
FROM enrol
JOIN courses on enrol.course_id = courses.id
JOIN roles ON enrol.role = roles.id
JOIN periods ON periods.id = courses.period
WHERE user_uid = $1
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Enrolment{}, err
	}
	var enrolment []Enrolment
	for rows.Next() {
		var id, year int
		var role, name, code, period string
		err := rows.Scan(&id, &role, &name, &code, &period, &year)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		enrolment = append(enrolment, Enrolment{
			Course: Course{
				ID:         id,
				Name:       name,
				CourseCode: code,
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

// GetAssessment will fetch the entire list of assessments for some student
func (store *Store) GetAssessment(uid string) ([]Assessment, error) {
	query := `
WITH statuses AS (
	SELECT
		submissions.assessment_id,
		type
	FROM submissions
	JOIN test_result_types ON test_result_types.id = submissions.result_id
	ORDER BY timestamp LIMIT 1)
SELECT
	assessment.id,
	courses.id,
	assessment.name,
	assessment_types.type,
	coalesce(statuses.type, 'untested')
FROM courses
JOIN enrol ON enrol.course_id = courses.id
JOIN assessment ON assessment.course_id = enrol.course_id
JOIN assessment_types ON assessment.type = assessment_types.id
FULL OUTER JOIN statuses ON statuses.assessment_id = assessment.id
WHERE enrol.user_uid = $1
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Assessment{}, err
	}
	var assessment []Assessment
	for rows.Next() {
		var id, courseID int
		var name, assessmentType, testResult string
		err := rows.Scan(
			&id,
			&courseID,
			&name,
			&assessmentType,
			&testResult,
		)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		assessment = append(assessment, Assessment{
			ID:         id,
			Name:       name,
			Type:       assessmentType,
			CourseID:   courseID,
			TestResult: testResult,
		})
	}
	rows.Close()
	return assessment, nil
}

// GetSubmissions will fetch all the submissions for a given user
func (store *Store) GetSubmissions(uid string) ([]Submission, error) {
	query := `
SELECT
	submissions.id,
	title,
	description,
	feedback,
	timestamp,
	name as assessment_name,
	course_id,
	assessment_id,
	test_result_types.type AS test_result,
	warnings,
	errors
FROM submissions
JOIN assessment ON assessment_id = assessment.id
JOIN test_results ON result_id = test_results.id
JOIN test_result_types ON test_results.type = test_result_types.id
WHERE uid = $1`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Submission{}, err
	}
	var submissions []Submission
	for rows.Next() {
		var title, assessmentName, description, feedback, testResult, warnings, errors string
		var id, courseID, assessmentID int
		var timestamp time.Time
		err := rows.Scan(
			&id,
			&title,
			&description,
			&feedback,
			&timestamp,
			&assessmentName,
			&courseID,
			&assessmentID,
			&testResult,
			&warnings,
			&errors,
		)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		submissions = append(submissions, Submission{
			ID:             id,
			AssessmentName: assessmentName,
			AssessmentID:   assessmentID,
			CourseID:       courseID,
			Title:          title,
			Description:    description,
			Feedback:       feedback,
			TestResult:     testResult,
			Warnings:       warnings,
			Errors:         errors,
			Timestamp:      timestamp,
		})
	}
	rows.Close()
	return submissions, nil
}

// GetTutorialEnrolment will fetch a list of tutorials for the
// given UID. If the user is enrolled in the course as a tutor,
// they are the tutorials the tutor is scheduled to mark.
func (store *Store) GetTutorialEnrolment(uid string) ([]TutorialEnrolment, error) {
	query := `
SELECT
	id, name
FROM tutorial_enrol
JOIN tutorials
ON tutorial_enrol.tutorial_id = tutorials.id
WHERE user_uid = $1
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []TutorialEnrolment{}, err
	}
	var enrolment []TutorialEnrolment
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		enrolment = append(enrolment, TutorialEnrolment{})
	}
	rows.Close()
	return enrolment, nil
}
