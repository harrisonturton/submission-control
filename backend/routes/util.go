package routes

import (
	"github.com/pkg/errors"
	"net/http"
)

const (
	errBadRequest          = "bad request"
	errUnauthorized        = "unauthorized"
	errNotFound            = "not found"
	errInternalServerError = "internal server error"
)

// queryURL will return the first value in a URL query string. It expects
// the parameter to have 1 value only, and will otherwise return an error.
func queryURL(key string, r *http.Request) (string, error) {
	values, ok := r.URL.Query()[key]
	if !ok || len(values) != 1 {
		return "", errors.New("failed to get parameter from query url")
	}
	return values[0], nil
}

// post is middleware that rejects requests that do not have a POST
// http method.
func post(handler http.HandlerFunc) http.HandlerFunc {
	return checkHTTPMethod(http.MethodPost, handler)
}

// get is middleware that rejects requests that do not have a POST
// http method.
func get(handler http.HandlerFunc) http.HandlerFunc {
	return checkHTTPMethod(http.MethodGet, handler)
}

// checkHTTPMethod is middleware that rejects requests that do not have
// the specified HTTP method.
func checkHTTPMethod(httpMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethod {
			writeNotFound(w)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// writeInternalServerError will write an internal server error to the
// ResponseWriter.
func writeInternalServerError(w http.ResponseWriter) {
	http.Error(w, errInternalServerError, http.StatusInternalServerError)
}

// writeBadRequest will write an bad request error message to the
// ResponseWriter.
func writeBadRequest(w http.ResponseWriter) {
	http.Error(w, errBadRequest, http.StatusBadRequest)
}

// writeUnauthorized will write an unauthorized error message to the
// ResponseWriter.
func writeUnauthorized(w http.ResponseWriter) {
	http.Error(w, errUnauthorized, http.StatusUnauthorized)
}

// writeNotFound will write a not found error message to the ResponseWriter.
func writeNotFound(w http.ResponseWriter) {
	http.Error(w, errNotFound, http.StatusNotFound)
}
