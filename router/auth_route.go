package router

import (
	"encoding/json"
	"net/http"
)

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
}
