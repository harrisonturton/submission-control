package store

// Account represents the login and account
// information that is required for every user.
type Account struct {
	Email        string
	Firstname    string
	Lastname     string
	PasswordHash string
	UID          string
}

// User represents a user with account information,
// and a set of roles.
type User struct {
	Account
	Roles []string
}
