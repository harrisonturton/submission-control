package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

// Store represents the database. It does NOT represent
// a connection, since database/sql manages a connection
// pool under the hood.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance, and ensures
// it can connect to the database.
func NewStore(db *sql.DB) (*Store, error) {
	err := db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the database")
	}
	return &Store{db}, nil
}

// GetAccountByEmail will fetch the account information for a user
// using their email address.
func (store *Store) GetAccountByEmail(email string) (*Account, error) {
	var name, password, uid string
	query := "SELECT name, password, uid FROM users WHERE email = $1"
	err := store.db.QueryRow(query, email).Scan(&name, &password, &uid)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch account")
	}
	return &Account{
		Name:     name,
		Email:    email,
		Password: password,
		UID:      uid,
	}, nil
}
