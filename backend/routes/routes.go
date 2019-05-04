package routes

import (
	"bytes"
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"io"
	"log"
	"net/http"
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
			return
		}
		resp, err := buildRefreshResponse(uid)
		if err != nil {
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

func studentUploadHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Println("Failed to get file")
			writeBadRequest(w)
			return
		}
		defer file.Close()
		io.Copy(&buf, file)
		contents := buf.String()
		w.Write([]byte(contents))
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

func needsAuthorization(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reject if not authorized
		if !request.IsAuthorized(r) {
			log.Println("Unauthorized")
			writeUnauthorized(w)
			return
		}
		// Else handle normally
		handler(w, r)
	})
}
