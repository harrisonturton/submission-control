package request

import (
	"context"
	"github.com/harrisonturton/submission-control/auth"
	"io/ioutil"
	"log"
	"net/http"
)

type key int

const (
	// AuthKey is the key that indexes the authentication information
	// in the requests context.
	AuthKey key = iota

	// BodyKey is the key that indexes the body data in the requests
	// context.
	BodyKey
)

// NewContextWithAuth verifies the JWT token (if it exists) in the
// request Authorization header. The result of this is stored in a boolean
// under the "auth" key in the new context.
func NewContextWithAuth(ctx context.Context, req *http.Request) context.Context {
	tokenHeader := req.Header["Authorization"]
	if len(tokenHeader) != 1 {
		return context.WithValue(ctx, AuthKey, false)
	}
	log.Printf("%v", tokenHeader[0])
	ok := auth.VerifyToken(tokenHeader[0])
	return context.WithValue(ctx, AuthKey, ok)
}

// NewContextWithBody reads the body of the request, and puts it in the context.
// We can only read the body once (since it is a ReadCloser), and so this lets
// multiple handlers/middleware read the body.
func NewContextWithBody(ctx context.Context, req *http.Request) context.Context {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return context.WithValue(ctx, BodyKey, "")
	}
	return context.WithValue(ctx, BodyKey, body)
}

// IsAuthorized will return true or false depending on whether the request was
// authenticated with the JWT token.
// It assumes the "auth" key has already been added to the context. This will
// crash if not true.
func IsAuthorized(r *http.Request) bool {
	return r.Context().Value(AuthKey).(bool)
}

// GetBody returns the body stored in the request context. It assumes this has
// already been added to the context.
// It assumes the "body" key has already been added to the context. This will
// crash if not true.
func GetBody(r *http.Request) []byte {
	return r.Context().Value(BodyKey).([]byte)
}
