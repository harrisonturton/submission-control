package routes

import (
	"encoding/json"
	"fmt"
	"github.com/harrisonturton/submission-control/backend/auth"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"net/http"
)

// These are the handlers called for each route, as specified in
// routes/router.go
//
// http.ServeMux expects a http.Handler, which is an interface for types
// with a ServeHTTP(http.ResponseWriter, *http.Request) function.
// Since this does not give any space to add extra dependencies (e.g.
// logger and database instances), I've manipulated closures to give the
// handler access to these.
//
// I could have implemented a custom Router{} type, with fields for various
// database and logger instances, but this is slightly cleaner and (hopefully)
// easier to test.

const (
	// These are the messages given for various error routes
	errorNotFoundMessage     = "not found"
	errorUnauthorizedMessage = "unauthorized"
)

// authHandler is called on the /auth route to request a new JWT token.
// It will authenticate the LoginRequest and generate a new token.
func authHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reject if not POST
		if r.Method != http.MethodPost {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Unmarshal POST data
		var login LoginRequest
		err := json.Unmarshal(request.GetBody(r), &login)
		if err != nil {
			badRequestHandler("unrecognised body").ServeHTTP(w, r)
			return
		}
		// Verify login information
		ok, err := auth.Authenticate(store, login.Email, login.Password)
		if !ok || err != nil {
			fmt.Println("Could not verify")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(login.Email)
		if err != nil {
			fmt.Println("Could not generate token")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := TokenResponse{
			StatusCode: http.StatusOK,
			Token:      token,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		log.Println(string(respBytes))
		w.Write(respBytes)
	})
}

// refreshHandler recieves a RefreshRequest. If the given token is valid,
// a new token is returned.
func refreshHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Reject is not POST
		if r.Method != http.MethodPost {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Unmarshal POST data
		var refresh RefreshRequest
		err := json.Unmarshal(request.GetBody(r), &refresh)
		if err != nil {
			badRequestHandler("unrecognised body").ServeHTTP(w, r)
			return
		}
		// Get claims from the token
		claims, err := auth.ParseToken(refresh.Token)
		if err != nil {
			log.Println("Failed to parse the token")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Generate new token
		token, err := auth.GenerateToken(claims.Email)
		if err != nil {
			log.Println("Failed to generate token")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := TokenResponse{
			StatusCode: http.StatusOK,
			Token:      token,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write(respBytes)
		log.Println(string(respBytes))
	})
}

func enrolHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Reject is not POST
		if r.Method != http.MethodGet {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Get user UID from the URL
		userUIDParams, ok := r.URL.Query()["uid"]
		if !ok || len(userUIDParams) != 1 {
			badRequestHandler("unrecognised UID").ServeHTTP(w, r)
			return
		}
		userUID := userUIDParams[0]
		// Get data
		courses, err := store.GetCoursesByUser(userUID)
		if err != nil {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := EnrolResponse{
			StatusCode: 200,
			Courses:    courses,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write(respBytes)
	})
}

func userHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			log.Println("Unauthorized")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Reject is not POST
		if r.Method != http.MethodGet {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Get user UID from the URL
		userEmailParams, ok := r.URL.Query()["email"]
		if !ok || len(userEmailParams) != 1 {
			log.Println("Invalid email")
			badRequestHandler("invalid email").ServeHTTP(w, r)
			return
		}
		userEmail := userEmailParams[0]
		// Get data
		user, err := store.GetUserByEmail(userEmail)
		if user == nil || err != nil {
			log.Printf("Could not find user: %v", err)
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := UserResponse{
			StatusCode: 200,
			User:       *user,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write(respBytes)
	})
}

func submissionHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			log.Println("Unauthorized")
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Reject is not POST
		if r.Method != http.MethodGet {
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Get user UID from the URL
		userUIDParams, ok := r.URL.Query()["uid"]
		if !ok || len(userUIDParams) != 1 {
			badRequestHandler("unrecognised UID").ServeHTTP(w, r)
			return
		}
		userUID := userUIDParams[0]
		submissions, err := store.GetSubmissionsForUser(userUID)
		if err != nil {
			log.Printf("Could not find submissions: %v", err)
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Send response
		resp := SubmissionResponse{
			StatusCode:  200,
			Submissions: submissions,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		w.Write(respBytes)
	})
}

func assessmentHandler(store *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsAuthorized(r) {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		// Get course ID from the URL
		userUIDParam, ok := r.URL.Query()["uid"]
		if !ok || len(userUIDParam) != 1 {
			log.Println("Invalid course ID")
			badRequestHandler("invalid user id").ServeHTTP(w, r)
			return
		}
		// Get assessments from the store
		userUID := userUIDParam[0]
		assessment, err := store.GetAssessmentForUser(userUID)
		if err != nil {
			log.Println("Could not find assessments")
			notFoundHandler().ServeHTTP(w, r)
			return
		}
		// Respond
		resp := AssessmentResponse{
			StatusCode: 200,
			Assessment: assessment,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			unauthorizedHandler().ServeHTTP(w, r)
			return
		}
		log.Println(string(respBytes))
		w.Write(respBytes)
	})
}

// notFoundHandler is called on other routes. It will return
// a 404 message.
func notFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusNotFound, errorNotFoundMessage)
		http.Error(w, resp, http.StatusNotFound)
		log.Println(resp)
	})
}

// unauthorizedHandler responds with a status unauthorized message
func unauthorizedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusUnauthorized, errorUnauthorizedMessage)
		http.Error(w, resp, http.StatusUnauthorized)
		log.Println(resp)
	})
}

// badRequestHandler responds with a bad request code and a custom message
func badRequestHandler(message string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf(
			`{"status":%d,"message":"%s"}`, http.StatusBadRequest, message)
		http.Error(w, resp, http.StatusBadRequest)
		log.Println(resp)
	})
}
