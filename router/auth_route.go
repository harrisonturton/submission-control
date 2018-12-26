package router

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var jwtTimeout = time.Minute * 5

var signingKey = []byte("asdgasdfglhjasbfjhkagsdlvgy gasdyfg")

// Claims is the data passed with the JWT for
// authentication
type Claims struct {
	Email string
	jwt.StandardClaims
}

/*
1) On client, ask for username and password
2) Exchange the username and password for a time-limited access token via HTTPS. Use jwt-go on the server
   to create the token. Use bcrypt to encrypt and compare passwords.
3) Add the recieved access token to the request header for any RESTful API requiring authorization
4) On the server, add an access token checker middleware for those routes. JWT tokens have an expire (exp)
   and not before (nbf) timestamp. JWT validates those when it parses the token from the header.
5) On client, periodically refresh the token. Our tokens expire in 5 minutes. I refresh them every 4 minutes
*/

func (router *Router) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		router.notFoundHandler(w, r)
		return
	}
	var login LoginRequest
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		router.logger.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "cannot decode body", http.StatusBadRequest)
		return
	}
	router.logger.Printf("Email: %s // Password: %s\n", login.Email, login.Password)
	account, err := router.store.GetAccountByEmail(login.Email)
	if err != nil {
		router.logger.Printf("Cannot find user: %v\n", err)
		http.Error(w, "request failed", http.StatusBadRequest)
		return
	}
	router.logger.Printf("Uid for %s is %s\n", account.Firstname, account.UID)
	result, err := checkPassword(login.Password, account.Password)
	if err != nil {
		router.logger.Printf("Failed to check password: %v\n", err)
		http.Error(w, "request failed", http.StatusBadRequest)
		return
	}
	if result {
		token, err := generateToken(account.Email)
		if err != nil {
			http.Error(w, "authentication failed", http.StatusBadRequest)
			return
		}
		router.logger.Printf("Successfully authenticated %s (%s)\n", account.Firstname, account.UID)
		fmt.Fprintf(w, token)
	} else {
		router.logger.Printf("Failed to authenticate %s (%s)\n", account.Firstname, account.UID)
		fmt.Fprintf(w, "Failed to authenticate...")
	}
}

func checkPassword(attempt string, hash string) (bool, error) {
	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(bytes, []byte(attempt))
	return err == nil, nil
}

func generateToken(email string) (string, error) {
	claims := Claims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtTimeout).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
