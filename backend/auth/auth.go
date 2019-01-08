package auth

import (
	"encoding/hex"
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
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	// SigningKey is the key used to sign each JWT token, proving
	// that it was given by this server. KEEP PRIVATE!
	SigningKey = "signing-key"

	// TokenTimeout is the duration a client can use the token to
	// access authenticated resources, after which it will be
	// rejected.
	TokenTimeout = time.Minute * 5
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
func GenerateToken(email string) (string, error) {
	claims := Claims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTimeout).Unix(),
			Issuer:    "server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SigningKey))
}

// Authenticate will verify a login attempt.
func Authenticate(store *store.Store, email string, passwordAttempt string) (bool, error) {
	// Get account data
	account, err := store.GetAccountByEmail(email)
	if err != nil {
		return false, err
	}
	// password hash is read as a hex string from the database (where it
	// is stored in binary), but bcrypt expects a []byte
	passwordHashBytes, err := hex.DecodeString(account.PasswordHash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(passwordHashBytes, []byte(passwordAttempt))
	return err == nil, err
}
