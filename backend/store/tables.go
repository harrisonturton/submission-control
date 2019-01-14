package store

import (
	"time"
)

// User represents the login and account
// information that is required for every user.
type User struct {
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	UID       string   `json:"uid"`
	Roles     []string `json:"roles"`

	PasswordHash string
}

// Course is the information about a single course.
// A course is uniquely identified by the course code
// and running semester (Period).
type Course struct {
	Name        string       `json:"name"`
	CourseCode  string       `json:"course_code"`
	Period      Semester     `json:"period"`
	Assignments []Assessment `json:"assignments"`
	Labs        []Assessment `json:"labs"`
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
	Name    string    `json:"name"`
	DueDate time.Time `json:"due_date"`
}
