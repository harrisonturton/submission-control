package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // The postgres driver
	"github.com/pkg/errors"
)

// Store represents the database. It does NOT represent
// a connection, since database/sql manages a connection
// pool under the hood.
type Store struct {
	db     *sql.DB
	config Config
}

// Config is the data required to connect to the postgres database
type Config struct {
	User    string
	DbName  string
	SslMode string
}

// ConnStr builds the connection string to pass to the postgres
// driver.
func (config *Config) ConnStr() string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=%s", config.User, config.DbName, config.SslMode)
}

// NewStore creates a new Store instance, and ensures
// it can connect to the database.
func NewStore(config Config) (*Store, error) {
	db, err := sql.Open("postgres", config.ConnStr())
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return &Store{db, config}, nil
}

// GetAccountByEmail will fetch the account information for a user
// using their email address.
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
