package routes

import (
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"net/http"
	"strconv"
)

func authHandler(store store.Reader) http.HandlerFunc {
	return post(func(w http.ResponseWriter, r *http.Request) {
		// Unmarshal the POST body
		var login LoginRequest
		err := json.Unmarshal(request.GetBody(r), &login)
		if err != nil {
			writeBadRequest(w)
			return
		}
		// Build the response
		resp, err := buildAuthResponse(store, login)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		w.Write(resp)
	})
}

func refreshHandler() http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			log.Println("Bad request to " + r.URL.Path)
			return
		}
		resp, err := buildRefreshResponse(uid)
		if err != nil {
			log.Println("Unauthorized access to " + r.URL.Path)
			writeUnauthorized(w)
			return
		}
		w.Write(resp)
	}))
}

func userHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildUserResponse(store, uid)
		if err != nil {
			log.Println("failed to build user response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func assessmentHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildAssessmentResponse(store, uid)
		if err != nil {
			log.Println("failed to build assessment response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func submissionsHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildSubmissionsResponse(store, uid)
		if err != nil {
			log.Println("failed to build submissions response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func studentUploadHandler(store *store.Store) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got student upload!")
		rawCourseID, err := queryURL("course_id", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		courseID, err := strconv.Atoi(rawCourseID)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildStudentUploadResponse(store, courseID, r.Body)
		if err != nil {
			log.Println("Failed to build studentUploadResponse " + err.Error())
			writeInternalServerError(w)
			return
		}
		log.Println("Got upload respones: " + string(resp))
		w.Write(resp)
	}))
}

func submissionUploadHandler(store *store.Store) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got submission upload!")
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		title, err := queryURL("title", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		description, err := queryURL("description", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		rawAssessmentID, err := queryURL("assessment_id", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		assessmentID, err := strconv.Atoi(rawAssessmentID)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildSubmissionUploadResponse(store, uid, title, description, assessmentID, r.Body)
		if err != nil {
			log.Println("Failed to build submissionUploadResponse " + err.Error())
			writeInternalServerError(w)
			return
		}
		log.Println("Got upload response: " + string(resp))
		w.Write(resp)
	}))
}

func logHandler(logger *log.Logger) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		// Read data somehow?
		// io.Copy(&buf, _)
		// log contents
		log.Println("Received remote logging")
	}))
}

func tutorialHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildTutorialResponse(store, uid)
		if err != nil {
			log.Println("failed to build tutorial response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}
