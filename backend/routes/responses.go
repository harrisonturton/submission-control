package routes

import (
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/auth"
	"github.com/harrisonturton/submission-control/backend/store"
	"github.com/pkg/errors"
)

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
		return nil, err
	}
	enrollments, err := store.GetEnrolment(uid)
	if err != nil {
		return nil, err
	}
	assessment, err := store.GetAssessment(uid)
	if err != nil {
		return nil, err
	}
	submissions, err := store.GetSubmissions(uid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(StudentStateResponse{
		User:        *user,
		Assessment:  assessment,
		Submissions: submissions,
		Enrolled:    enrollments,
	})
}
