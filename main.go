package main

import (
	"database/sql"
	"fmt"
	dbStore "github.com/harrisonturton/submission-control/db"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	dbUser    = "harrisonturton"
	dbName    = "submission_control"
	dbSslMode = "disable"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	store, err := createDb(dbUser, dbName, dbSslMode)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()
	// Now do some work with the database
	user, err := store.GetUser("u6386433")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("UID: %s, Name: %s, Roles: %v\n", user.UID, user.Name, user.Roles)
	logger.Println("Done!")
}

func createDb(dbUser, dbName, dbSslMode string) (*dbStore.Store, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", dbUser, dbName, dbSslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return dbStore.NewStore(db)
}
