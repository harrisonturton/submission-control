package store

import (
	"database/sql"
	"github.com/pkg/errors"
)

// Store represents the database. It does NOT
// represent a single connection, since database/sql
// manages a connection pool beneath the hood.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance, and ensures
// the database can be connected to.
func NewStore(db *sql.DB) (*Store, error) {
	err := db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return &Store{db}, nil
}

// GetUserByEmail will return a single user with a
// matching email address.
func (store *Store) GetUserByEmail(email string) (*User, error) {
	var firstname, lastname, uid string
	var passwordHash []byte
	query := "SELECT first_name, last_name, password, uid FROM users WHERE email = $1"
	err := store.db.QueryRow(query, email).Scan(&firstname, &lastname, &passwordHash, &uid)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch account")
	}
	return &User{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		PasswordHash: string(passwordHash),
		UID:          uid,
	}, nil
}
