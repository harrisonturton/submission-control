package routes

import (
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
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

func stateHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildStudentStateResponse(store, uid)
		if err != nil {
			log.Println("failed to build state response")
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
