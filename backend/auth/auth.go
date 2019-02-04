package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/harrisonturton/submission-control/backend/store"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// Claims is the data passed with the JWT for
// authentication
type Claims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

var (
	// SigningKey is the key used to sign each JWT token, proving
	// that it was given by this server. KEEP PRIVATE!
	SigningKey = "signing-key"

	// TokenTimeout is the duration a client can use the token to
	// access authenticated resources, after which it will be
	// rejected.
	TokenTimeout = time.Minute * 10
)

// ParseToken parses a token into a Claims instance
func ParseToken(rawToken string) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SigningKey), nil
	})
	return &claims, err
}

// VerifyToken will try and verify the JWT token.
func VerifyToken(rawToken string) bool {
	var claims Claims
	token, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}
	_, ok := token.Claims.(*Claims)
	return ok && token.Valid
}

// GenerateToken will generate a new JWT token from an email address.
func GenerateToken(uid string) (string, error) {
	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTimeout).Unix(),
			Issuer:    "server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SigningKey))
}

// Authenticate will verify a login attempt.
func Authenticate(store store.Reader, uid string, passwordAttempt string) (bool, error) {
	// Get user data
	user, err := store.GetUser(uid)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordAttempt))
	return err == nil, err
}
