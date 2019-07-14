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
WITH
	test_results AS (
		SELECT
			submissions.title,
			submissions.assessment_id,
			test_result_types.type
		FROM submissions
		JOIN test_results ON submissions.id = test_results.submission_id
		JOIN test_result_types ON test_result_types.id = test_results.type),
	enrolled_courses AS (
		SELECT
			*
		FROM courses
		JOIN enrol ON enrol.course_id = courses.id
		WHERE enrol.uid=$1)
SELECT
	assessment.id,
	enrolled_courses.course_id,
	assessment.name,
	assessment_types.type,
	coalesce(test_results.type, 'untested')
FROM assessment
JOIN enrolled_courses ON assessment.course_id = enrolled_courses.course_id
JOIN assessment_types ON assessment_types.id = assessment.type
LEFT JOIN test_results ON assessment.id = test_results.assessment_id`
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
	uid,
	title,
	description,
	feedback,
	timestamp,
	name AS assessment_name,
	course_id,
	assessment_id,
	COALESCE(test_result_types.type, 'untested'),
	warnings,
	errors
FROM submissions
JOIN assessment on assessment.id = submissions.assessment_id
LEFT JOIN test_results ON test_results.submission_id = submissions.id
LEFT JOIN test_result_types ON test_results.type = test_result_types.id
WHERE uid=$1
`
	rows, err := store.db.Query(query, uid)
	if err != nil {
		return []Submission{}, err
	}
	var submissions []Submission
	for rows.Next() {
		var warnings, errors *string
		var title, uid, assessmentName, description, feedback, testResult string
		var id, courseID, assessmentID int
		var timestamp time.Time
		err := rows.Scan(
			&id,
			&uid,
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
			UID:            uid,
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

// GetCourse will fetch a singular course.
func (store *Store) GetCourse(courseID int) (*Course, error) {
	query := `
SELECT
	courses.id,
	courses.code,
	courses.name,
	courses.year,
	periods.period
FROM courses
JOIN periods ON periods.id = courses.period
WHERE courses.id = $1
`
	var id, year int
	var name, courseCode, period string
	err := store.db.
		QueryRow(query, courseID).
		Scan(&id, &courseCode, &name, &year, &period)
	if err != nil {
		return nil, err
	}
	return &Course{
		ID:         id,
		Name:       name,
		CourseCode: courseCode,
		Period:     period,
		Year:       year,
	}, nil
}

// GetTutorial will find a tutorial with the given ID
func (store *Store) GetTutorial(tutorialID int) (*Tutorial, error) {
	query := `SELECT * FROM tutorials WHERE tutorials.id = $1`
	var id, courseID int
	var name string
	err := store.db.
		QueryRow(query, tutorialID).
		Scan(&id, &name, &courseID)
	if err != nil {
		return nil, err
	}
	tutors, err := store.getTutorsForTutorial(tutorialID)
	if err != nil {
		return nil, err
	}
	students, err := store.getStudentsForTutorial(tutorialID)
	if err != nil {
		return nil, err
	}
	submissions, err := store.getSubmissionsForTutorial(tutorialID)
	if err != nil {
		return nil, err
	}
	assessment, err := store.getAssessmentForTutorial(tutorialID)
	if err != nil {
		return nil, err
	}
	return &Tutorial{
		ID:          id,
		Name:        name,
		Tutors:      tutors,
		Students:    students,
		Submissions: submissions,
		Assessment:  assessment,
	}, nil
}

func (store *Store) getTutorsForTutorial(tutorialID int) ([]User, error) {
	query := ` 
SELECT
	users.first_name,
	users.last_name,
	users.email,
	users.uid
FROM tutorial_enrol
JOIN tutorials ON tutorial_enrol.tutorial_id = tutorials.id
JOIN enrol ON tutorials.course_id = enrol.course_id AND tutorial_enrol.uid = enrol.uid
JOIN roles ON enrol.role = roles.id
JOIN users ON users.uid = tutorial_enrol.uid
WHERE roles.role = 'tutor' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []User{}, err
	}
	var tutors []User
	for rows.Next() {
		var firstName, lastName, email, uid string
		err := rows.Scan(&firstName, &lastName, &email, &uid)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		tutors = append(tutors, User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			UID:       uid,
			Enrolment: []Enrolment{},
		})
	}
	rows.Close()
	if len(tutors) == 0 {
		return []User{}, nil
	}
	return tutors, nil
}

func (store *Store) getStudentsForTutorial(tutorialID int) ([]User, error) {
	query := `
SELECT
	users.uid,
	users.first_name,
	users.last_name,
	users.email,
	users.password
FROM tutorial_enrol
JOIN tutorials ON tutorial_enrol.tutorial_id = tutorials.id
JOIN enrol ON tutorials.course_id = enrol.course_id AND tutorial_enrol.uid = enrol.uid
JOIN roles ON enrol.role = roles.id
JOIN users ON users.uid = tutorial_enrol.uid
WHERE roles.role = 'student' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []User{}, err
	}
	var students []User
	for rows.Next() {
		var uid, firstName, lastName, email string
		var passwordHash []byte
		err := rows.Scan(&uid, &firstName, &lastName, &email, &passwordHash)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		students = append(students, User{
			Email:        email,
			FirstName:    firstName,
			LastName:     lastName,
			UID:          uid,
			PasswordHash: string(passwordHash),
			Enrolment:    []Enrolment{},
		})
	}
	rows.Close()
	if len(students) == 0 {
		return []User{}, nil
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
JOIN test_results ON submissions.id = test_results.submission_id
JOIN test_result_types ON test_results.type = test_result_types.id
WHERE roles.role = 'student' AND tutorials.id = $1`
	rows, err := store.db.Query(query, tutorialID)
	if err != nil {
		return []Submission{}, err
	}
	var submissions []Submission = []Submission{}
	for rows.Next() {
		var warnings, errors *string
		var title, assessmentName, description, feedback, testResult string
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
