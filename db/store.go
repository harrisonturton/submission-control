package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

var (
	errNotFound = errors.New("database item does not exist")
)

// Store represents the database. Each instance
// does NOT represent a connection, since database/sql
// manages a connection pool under the hood.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance, and pings
// the database to ensure a connection is working.
func NewStore(db *sql.DB) (*Store, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to the database")
	}
	return &Store{db}, nil
}

// GetUser will fetch a user from the database with the given UID.
func (store *Store) GetUser(uid string) (*User, error) {
	// Get the users name
	name, err := store.fetchUserData(uid)
	if err != nil {
		return nil, err
	}
	// Get the users roles
	roles, err := store.fetchRoles(uid)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:  name,
		UID:   uid,
		Roles: roles,
	}, nil
}

// fetchRoles will get an array of roles for a user
func (store *Store) fetchRoles(uid string) ([]string, error) {
	var roles []string
	rows, err := store.db.Query("SELECT role FROM roles WHERE uni_id = $1", uid)
	if err == sql.ErrNoRows {
		return roles, errNotFound
	}
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			continue
		}
		roles = append(roles, role)
	}
	rows.Close()
	return roles, nil
}

// fetchUserData will get data for a single row in the user table
func (store *Store) fetchUserData(uid string) (string, error) {
	var name string
	err := store.db.QueryRow("SELECT name FROM users WHERE uni_id = $1", uid).Scan(&name)
	if err == sql.ErrNoRows {
		return "", errNotFound
	}
	if err != nil {
		return "", err
	}
	return name, nil
}

// Close will close the connection to the database
func (store *Store) Close() error {
	return store.db.Close()
}
