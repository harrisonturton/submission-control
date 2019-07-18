package routes

import (
	"errors"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
)

// StudentUploadItem is the row of the .csv
// that is uploaded to populate the student database.
type StudentUploadItem struct {
	Student store.User
}

// Expects format { uid, firstname, lastname, email }
func parseStudentUpload(data [][]string) ([]StudentUploadItem, error) {
	results := []StudentUploadItem{}
	log.Println("Inside parseStudentUpload")
	log.Printf("Data: %v\nWith length: %d\n", data, len(data))
	for _, row := range data[1:] {
		log.Printf("Working with row %v\n", row)
		if len(row) != 4 {
			log.Printf("Row length is not 6 but %d\n", len(row))
			continue
		}
		item := StudentUploadItem{
			Student: store.User{
				UID:       row[0],
				FirstName: row[1],
				LastName:  row[2],
				Enrolment: []store.Enrolment{},
				Email:     row[3],
			},
		}
		results = append(results, item)
	}
	if len(results) == 0 {
		return nil, errors.New("table is either empty, or invalid")
	}
	return results, nil
}
