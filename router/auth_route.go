package router

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

func (router *Router) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		router.notFoundHandler(w, r)
		return
	}
	// Decode login request
	login, err := decodeLoginRequest(r.Body)
	if err != nil {
		router.logger.Printf("Error decoding login request: %v\n", err)
		http.Error(w, "cannot decode body", http.StatusBadRequest)
	}
	// Get account from email
	account, err := router.store.GetAccountByEmail(login.Email)
	if err != nil {
		router.logger.Printf("Cannot find user for email %s: %v\n", login.Email, err)
		http.Error(w, "request failed", http.StatusBadRequest)
		return
	}
	// Verify login information
	ok, err := checkPassword(login.Password, account.Password)
	if !ok || err != nil {
		router.logger.Printf("Failed to check password: %v\n", err)
		http.Error(w, "request failed", http.StatusBadRequest)
		return
	}
	// Generate new token
	token, err := generateToken(account.Email)
	if err != nil {
		router.logger.Printf("Error generating token for %s: %v\n", account.Firstname, err)
		http.Error(w, "request failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, token)
}

func generateToken(email string) (string, error) {
	claims := Claims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTimeout).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SigningKey)
}

func decodeLoginRequest(body io.Reader) (*LoginRequest, error) {
	var login LoginRequest
	err := json.NewDecoder(body).Decode(&login)
	return &login, err
}

func checkPassword(plaintextAttempt string, actualHash string) (bool, error) {
	bytes, err := hex.DecodeString(actualHash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(bytes, []byte(plaintextAttempt))
	return err == nil, nil
}
