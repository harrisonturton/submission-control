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
	for _, row := range data {
		if len(row) != 6 {
			continue
		}
		tutorialID, err := strconv.Atoi(row[3])
		if err != nil {
			log.Println("Invalid tutorial ID: " + row[3])
			continue
		}
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
