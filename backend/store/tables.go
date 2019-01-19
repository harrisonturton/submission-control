package store

// User represents the login and account
// information that is required for every user.
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UID       string `json:"uid"`

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

// Semester represents the time when a course is being run
type Semester int

// These are the different periods in which a course can be run
const (
	Summer Semester = iota
	First
	Autumn
	Winter
	Second
	Spring
)

// Assessment represents a single piece of assessment, whether
// it be a proper assignment or something for a lab.
type Assessment struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	CourseID int    `json:"course_id"`
}

// Submission is a single submission made by a user to an assessment item.
type Submission struct {
	ID           int    `json:"id"`
	AssessmentID int    `json:"assessment_id"`
	CourseID     int    `json:"course_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Feedback     string `json:"feedback"`
}
