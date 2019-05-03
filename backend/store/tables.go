package store

import "time"

// User represents the login and account
// information that is required for every user.
type User struct {
	Email     string      `json:"email"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	UID       string      `json:"uid"`
	Enrolment []Enrolment `json:"enrolment"`

	PasswordHash string
}

// Course is the information about a single course.
// A course is uniquely identified by the course code
// and running semester (Period).
type Course struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CourseCode string `json:"course_code"`
	Period     string `json:"period"`
	Year       int    `json:"year"`
}

// Tutorial represents a lab/tutorial space for a course
type Tutorial struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Tutors      []string     `json:"tutors"`
	Students    []User       `json:"students"`
	Submissions []Submission `json:"submissions"`
	Assessment  []Assessment `json:"assessment"`
}

// TutorialEnrolment shows the tutorials a user
// is enrolled in.
type TutorialEnrolment struct {
	Tutorial Tutorial `json:"tutorial"`
	Course   Course   `json:"course"`
}

// These are the different periods in which a course
// can be run, each year
const (
	PeriodFirst  = "first"
	PeriodSecond = "second"
	PeriodSummer = "summer"
	PeriodAutumn = "autumn"
	PeriodWinter = "winter"
	PeriodSpring = "spring"
)

// Assessment represents a single piece of assessment, whether
// it be a proper assignment or something for a lab.
type Assessment struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	CourseID   int    `json:"course_id"`
	TestResult string `json:"test_result"`
}

// Submission is a single submission made by a user to an assessment item.
type Submission struct {
	ID             int    `json:"id"`
	UID            string `json:"uid"`
	AssessmentName string `json:"assessment_name"`
	AssessmentID   int    `json:"assessment_id"`
	CourseID       int    `json:"course_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Feedback       string `json:"feedback"`

	TestResult string    `json:"test_result"`
	Warnings   string    `json:"warnings"`
	Errors     string    `json:"errors"`
	Timestamp  time.Time `json:"timestamp"`
}

// Enrolment maps a user to the course they are, or have previously been,
// enrolled in.
type Enrolment struct {
	Course Course `json:"course"`
	UID    string `json:"uid"`
	Role   string `json:"role"`
}
