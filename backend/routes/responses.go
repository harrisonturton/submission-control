package routes

import (
	"encoding/csv"
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/auth"
	"github.com/harrisonturton/submission-control/backend/store"
	"github.com/pkg/errors"
	"io"
	"log"
)

// TokenResponse responds with a JWT token.
type TokenResponse struct {
	Token string `json:"token"`
}

// StudentStateResponse contains all the state data required to render
// a students page on the client.
type StudentStateResponse struct {
	User        store.User         `json:"user"`
	Assessment  []store.Assessment `json:"assessment"`
	Submissions []store.Submission `json:"submissions"`
	Enrolled    []store.Enrolment  `json:"enrolled"`
}

// UserResponse contains data for a single user and their
// enrolled courses.
type UserResponse struct {
	User      store.User        `json:"user"`
	Enrolment []store.Enrolment `json:"enrolment"`
}

// AssessmentResponse contains a list of assessments,
// for every course a user is enrolled in. Served on the
// /assessment GET endpoint.
type AssessmentResponse struct {
	Assessment []store.Assessment `json:"assessment"`
}

// SubmissionsResponse contains all the submissions the user
// has made.
type SubmissionsResponse struct {
	Submissions []store.Submission `json:"submissions"`
}

// TutorialResponse is given on the /tutorials endpoint
type TutorialResponse struct {
	Tutorials []store.TutorialEnrolment `json:"tutorials"`
}

// StudentRecord is the type for each row in a .csv
type StudentRecord struct {
	Firstname string
	Lastname  string
	Tutorials []string
}

func buildAuthResponse(store store.Reader, login LoginRequest) ([]byte, error) {
	ok, err := auth.Authenticate(store, login.UID, login.Password)
	if !ok || err != nil {
		return nil, errors.New("unauthorized")
	}
	token, err := auth.GenerateToken(login.UID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return json.Marshal(TokenResponse{
		Token: token,
	})
}

func buildRefreshResponse(token string) ([]byte, error) {
	claims, err := auth.ParseToken(token)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	newToken, err := auth.GenerateToken(claims.UID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return json.Marshal(TokenResponse{
		Token: newToken,
	})
}

func buildStudentStateResponse(store store.Reader, uid string) ([]byte, error) {
	user, err := store.GetUser(uid)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		return nil, err
	}
	enrollments, err := store.GetEnrolment(uid)
	if err != nil {
		log.Printf("Error getting enrolments: %v\n", err)
		return nil, err
	}
	assessment, err := store.GetAssessment(uid)
	if err != nil {
		log.Printf("Error getting assessments: %v\n", err)
		return nil, err
	}
	submissions, err := store.GetSubmissions(uid)
	if err != nil {
		log.Printf("Error getting submissions: %v\n", err)
		return nil, err
	}
	return json.Marshal(StudentStateResponse{
		User:        *user,
		Assessment:  assessment,
		Submissions: submissions,
		Enrolled:    enrollments,
	})
}

func buildUserResponse(store store.Reader, uid string) ([]byte, error) {
	user, err := store.GetUser(uid)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		return nil, err
	}
	enrollment, err := store.GetEnrolment(uid)
	if err != nil {
		log.Printf("Error getting enrolment: %v\n", err)
		return nil, err
	}
	return json.Marshal(UserResponse{
		User:      *user,
		Enrolment: enrollment,
	})
}

func buildAssessmentResponse(store store.Reader, uid string) ([]byte, error) {
	assessment, err := store.GetAssessment(uid)
	if err != nil {
		log.Printf("Error getting assessment: %v\n", err)
		return nil, err
	}
	return json.Marshal(AssessmentResponse{
		Assessment: assessment,
	})
}

func buildSubmissionsResponse(store store.Reader, uid string) ([]byte, error) {
	submissions, err := store.GetSubmissions(uid)
	if err != nil {
		log.Printf("Error getting submissions: %v\n", err)
		return nil, err
	}
	return json.Marshal(SubmissionsResponse{
		Submissions: submissions,
	})
}

func buildStudentUploadResponse(store *store.Reader, data io.Reader) ([]byte, error) {
	r := csv.NewReader(data)
	_, err := r.ReadAll() // [][]string
	if err != nil {
		log.Println("Failed to read .csv form data")
		return nil, err
	}
	return nil, err
}

func buildTutorialResponse(store store.Reader, uid string) ([]byte, error) {
	tutorialEnrolment, err := store.GetTutorialEnrolment(uid)
	if err != nil {
		log.Printf("Error getting tutorials: %v\n", err)
		return nil, err
	}
	b, _ := json.Marshal(tutorialEnrolment[0])
	log.Print(string(b))
	return json.Marshal(TutorialResponse{
		Tutorials: tutorialEnrolment,
	})
}
