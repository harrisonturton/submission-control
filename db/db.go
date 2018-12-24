package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	conn := "user=harrisonturton dbname=submission_control sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done!")
	rows, err := db.Query("SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found row containing %q", s)
	}
	rows.Close()
}
