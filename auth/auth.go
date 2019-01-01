package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// Claims is the data passed with the JWT for
// authentication
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// SigningKey ...
const SigningKey = "signing-key"

// TokenTimeout ...
var TokenTimeout = time.Minute * 5

// Authenticate will try and verify the JWT token.
func Authenticate(rawToken string) bool {
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
	c, ok := token.Claims.(*Claims)
	log.Println(c.Email)
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
