package routes

import (
	"errors"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"strconv"
)

// StudentUploadItem is the row of the .csv
// that is uploaded to populate the student database.
type StudentUploadItem struct {
	Student    store.User
	TutorialID int
}

func parseStudentUpload(data [][]string) ([]StudentUploadItem, error) {
	results := []StudentUploadItem{}
	log.Println("Inside parseStudentUpload")
	log.Printf("Data:  %v\nWith length: %d\n", data, len(data))
	for _, row := range data {
		log.Printf("Working with row %v\n", row)
		if len(row) != 6 {
			log.Printf("Row length is not 6 but %d\n", len(row))
			continue
		}
		tutorialID, err := strconv.Atoi(row[3])
		if err != nil {
			log.Println("Invalid tutorial ID: " + row[3])
			continue
		}
		log.Println("With name: " + row[0] + " " + row[1])
		item := StudentUploadItem{
			Student: store.User{
				UID:          row[0],
				FirstName:    row[1],
				LastName:     row[2],
				Enrolment:    []store.Enrolment{},
				Email:        row[4],
				PasswordHash: row[5],
			},
			TutorialID: tutorialID,
		}
		results = append(results, item)
	}
	if len(results) == 0 {
		return nil, errors.New("table is either empty, or invalid")
	}
	return results, nil
}
