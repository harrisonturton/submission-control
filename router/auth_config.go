package router

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims is the data passed with the JWT for
// authentication
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	// SigningKey is the key used to sign & verify each token
	SigningKey = []byte("submission-app-key")

	// TokenTimeout is the lifetime of the time-limited token. If
	// a client tries to authenticate with an old token, it is rejected.
	TokenTimeout = time.Minute * 5
)
