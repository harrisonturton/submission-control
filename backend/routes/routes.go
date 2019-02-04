package routes

import (
	"encoding/json"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

const (
	errBadRequest          = "bad request"
	errUnauthorized        = "unauthorized"
	errNotFound            = "not found"
	errInternalServerError = "internal server error"
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

func queryURL(key string, r *http.Request) (string, error) {
	values, ok := r.URL.Query()[key]
	if !ok || len(values) != 1 {
		return "", errors.New("failed to get parameter from query url")
	}
	return values[0], nil
}

func post(handler http.HandlerFunc) http.HandlerFunc {
	return checkHTTPMethod(http.MethodPost, handler)
}

func get(handler http.HandlerFunc) http.HandlerFunc {
	return checkHTTPMethod(http.MethodGet, handler)
}

func checkHTTPMethod(httpMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethod {
			writeNotFound(w)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func writeInternalServerError(w http.ResponseWriter) {
	http.Error(w, errInternalServerError, http.StatusInternalServerError)
}

func writeBadRequest(w http.ResponseWriter) {
	http.Error(w, errBadRequest, http.StatusBadRequest)
}

func writeUnauthorized(w http.ResponseWriter) {
	http.Error(w, errUnauthorized, http.StatusUnauthorized)
}

func writeNotFound(w http.ResponseWriter) {
	http.Error(w, errNotFound, http.StatusNotFound)
}
