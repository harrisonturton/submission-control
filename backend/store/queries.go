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
WHERE uid = $1
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
WHERE enrol.uid = $1
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

// GetTutorialEnrolment will fetch all the tutorials for the courses
// in which the user is enrolled as a tutor. It will fetch all the tutorials,
// not just the ones assigned to that user.
func (store *Store) GetTutorialEnrolment(uid string) ([]TutorialEnrolment, error) {
	query := `
SELECT
  tutorials.id,
  tutorials.name,
  courses.id,
  courses.name,
  courses.code,
  periods.period,
  courses.year
FROM enrol
JOIN courses ON enrol.course_id = courses.id
JOIN tutorials ON tutorials.course_id = courses.id
JOIN periods ON courses.period = periods.id
WHERE enrol.uid = $1`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []TutorialEnrolment{}, err
	}
	var enrolment []TutorialEnrolment
	for rows.Next() {
		var tutorialID, courseID, courseYear int
		var tutorialName, courseName, courseCode, coursePeriod string
		err = rows.Scan(
			&tutorialID,
			&tutorialName,
			&courseID,
			&courseName,
			&courseCode,
			&coursePeriod,
			&courseYear,
		)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		tutors, err := store.getTutorsForTutorial(tutorialID)
		if err != nil {
			log.Println("Failed to get tutors: " + err.Error())
			continue
		}
		students, err := store.getStudentsForTutorial(tutorialID)
		if err != nil {
			log.Println("Failed to get students: " + err.Error())
			continue
		}
		submissions, err := store.getSubmissionsForTutorial(tutorialID)
		if err != nil {
			log.Println("Failed to get submissions: " + err.Error())
			continue
		}
		assessment, err := store.getAssessmentForTutorial(tutorialID)
		if err != nil {
			log.Println("Failed to get assessment: " + err.Error())
			continue
		}
		enrolment = append(enrolment, TutorialEnrolment{
			Tutorial: Tutorial{
				ID:          tutorialID,
				Name:        tutorialName,
				Students:    students,
				Tutors:      tutors,
				Submissions: submissions,
				Assessment:  assessment,
			},
			Course: Course{
				ID:         courseID,
				Name:       courseName,
				CourseCode: courseCode,
				Period:     coursePeriod,
				Year:       courseYear,
			},
		})
	}
	rows.Close()
	return enrolment, nil
}

func (store *Store) getTutorsForTutorial(tutorialID int) ([]string, error) {
	query := ` 
SELECT
  tutorial_enrol.uid
FROM tutorial_enrol
JOIN tutorials ON tutorial_enrol.tutorial_id = tutorials.id
JOIN enrol ON tutorials.course_id = enrol.course_id AND tutorial_enrol.uid = enrol.uid
JOIN roles ON enrol.role = roles.id
WHERE roles.role = 'tutor' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []string{}, err
	}
	var tutors []string
	for rows.Next() {
		var uid string
		err := rows.Scan(&uid)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		tutors = append(tutors, uid)
	}
	rows.Close()
	if len(tutors) == 0 {
		return []string{}, nil
	}
	return tutors, nil
}

func (store *Store) getStudentsForTutorial(tutorialID int) ([]string, error) {
	query := ` 
SELECT
  tutorial_enrol.uid
FROM tutorial_enrol
JOIN tutorials ON tutorial_enrol.tutorial_id = tutorials.id
JOIN enrol ON tutorials.course_id = enrol.course_id AND tutorial_enrol.uid = enrol.uid
JOIN roles ON enrol.role = roles.id
WHERE roles.role = 'student' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []string{}, err
	}
	var students []string
	for rows.Next() {
		var uid string
		err := rows.Scan(&uid)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		students = append(students, uid)
	}
	rows.Close()
	if len(students) == 0 {
		return []string{}, nil
	}
	return students, nil
}

// getSubmissionsForTutorial will fetch all the submissions for students in a given
// tutorial
func (store *Store) getSubmissionsForTutorial(tutorialID int) ([]Submission, error) {
	query := `
SELECT
  submissions.id,
  assessment.name,
  assessment.id,
  assessment.course_id,
  submissions.title,
  submissions.description,
  submissions.feedback,
  test_result_types.type,
  test_results.warnings,
  test_results.errors,
  submissions.timestamp
FROM tutorial_enrol
JOIN tutorials ON tutorial_enrol.tutorial_id = tutorials.id
JOIN enrol ON tutorials.course_id = enrol.course_id AND tutorial_enrol.uid = enrol.uid
JOIN roles ON enrol.role = roles.id
JOIN submissions ON submissions.uid = enrol.uid
JOIN assessment ON submissions.assessment_id = assessment.id AND enrol.course_id = assessment.course_id
JOIN test_results ON submissions.result_id = test_results.id
JOIN test_result_types ON test_results.type = test_result_types.id
WHERE roles.role = 'student' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []Submission{}, err
	}
	var submissions []Submission = []Submission{}
	for rows.Next() {
		var title, assessmentName, description, feedback, testResult, warnings, errors string
		var id, courseID, assessmentID int
		var timestamp time.Time
		err := rows.Scan(
			&id,
			&assessmentName,
			&assessmentID,
			&courseID,
			&title,
			&description,
			&feedback,
			&testResult,
			&warnings,
			&errors,
			&timestamp,
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

func (store *Store) getAssessmentForTutorial(tutorialID int) ([]Assessment, error) {
	query := `
SELECT
	assessment.id,
	assessment.course_id,
	assessment.name,
	assessment_types.type
FROM assessment
JOIN tutorials ON tutorials.id = assessment.course_id
JOIN assessment_types ON assessment_types.id = assessment.type
WHERE tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []Assessment{}, err
	}
	var assessment []Assessment = []Assessment{}
	for rows.Next() {
		var assessmentID, courseID int
		var assessmentName, assessmentType string
		err := rows.Scan(
			&assessmentID,
			&courseID,
			&assessmentName,
			&assessmentType,
		)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		assessment = append(assessment, Assessment{
			ID:         assessmentID,
			Name:       assessmentName,
			Type:       assessmentType,
			CourseID:   courseID,
			TestResult: "untested",
		})
	}
	rows.Close()
	return assessment, nil
}
