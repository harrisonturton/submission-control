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
SELECT
	assessment.id,
	assessment.course_id,
	name,
	assessment_types.type,
	test_result_types.type
FROM assessment
JOIN submissions ON submissions.assessment_id = assessment.id
JOIN assessment_types ON assessment_types.id = assessment.type
JOIN test_results ON submissions.result_id = test_results.id
JOIN test_result_types ON test_result_types.id = test_results.id
WHERE uid = $1
ORDER BY submissions.timestamp DESC
LIMIT 1
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
