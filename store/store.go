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

// GetAccountByEmail will return a single account with a
// matching email address.
func (store *Store) GetAccountByEmail(email string) (*Account, error) {
	var firstname, lastname, password, uid string
	query := "SELECT firstname, lastname, password, uid FROM users WHERE email = $1"
	err := store.db.QueryRow(query, email).Scan(&firstname, &lastname, &password, &uid)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch account")
	}
	return &Account{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
		UID:       uid,
	}, nil
}
