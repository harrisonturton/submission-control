package router

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/harrisonturton/submission-control/auth"
	"github.com/harrisonturton/submission-control/request"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (router *Router) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		router.notFoundHandler(w, r)
		return
	}
	// Decode login request
	login, err := decodeLoginRequest(request.GetBody(r))
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
	token, err := auth.GenerateToken(account.Email)
	if err != nil {
		router.logger.Printf("Error generating token for %s: %v\n", account.Firstname, err)
		http.Error(w, "request failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, token+"\n")
}

func decodeLoginRequest(body []byte) (*LoginRequest, error) {
	var login LoginRequest
	err := json.Unmarshal(body, &login)
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
