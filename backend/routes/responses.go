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

func buildRefreshResponse(uid string) ([]byte, error) {
	token, err := auth.GenerateToken(uid)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return json.Marshal(TokenResponse{
		Token: token,
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

func buildStudentUploadResponse(store *store.Store, data io.Reader) ([]byte, error) {
	r := csv.NewReader(data)
	rawTable, err := r.ReadAll()
	if err != nil {
		log.Println("Failed to read .csv form data")
		return nil, err
	}
	table, err := parseStudentUpload(rawTable)
	if err != nil {
		log.Println("Failed to parse .csv form data")
		return nil, err
	}
	for _, row := range table {
		err := store.WriteUser(row.Student)
		if err != nil {
			log.Println("Failed to write user")
			continue
		}
		err = store.WriteTutorialEnrolment(row.Student.UID, row.TutorialID)
		if err != nil {
			log.Println("Failed to write enrolment for " + row.Student.UID)
			continue
		}
	}
	return nil, nil
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
