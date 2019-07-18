package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/harrisonturton/submission-control/backend/ci"
	"github.com/harrisonturton/submission-control/backend/request"
	"github.com/harrisonturton/submission-control/backend/store"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	mb = 1 << 20
)

func authHandler(store store.Reader) http.HandlerFunc {
	return post(func(w http.ResponseWriter, r *http.Request) {
		// Unmarshal the POST body
		var login LoginRequest
		err := json.Unmarshal(request.GetBody(r), &login)
		if err != nil {
			writeBadRequest(w)
			return
		}
		// Build the response
		resp, err := buildAuthResponse(store, login)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		w.Write(resp)
	})
}

func refreshHandler() http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			log.Println("Bad request to " + r.URL.Path)
			return
		}
		resp, err := buildRefreshResponse(uid)
		if err != nil {
			log.Println("Unauthorized access to " + r.URL.Path)
			writeUnauthorized(w)
			return
		}
		w.Write(resp)
	}))
}

func userHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildUserResponse(store, uid)
		if err != nil {
			log.Println("failed to build user response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func assessmentHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildAssessmentResponse(store, uid)
		if err != nil {
			log.Println("failed to build assessment response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func submissionsHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildSubmissionsResponse(store, uid)
		if err != nil {
			log.Println("failed to build submissions response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func submissionFileHandler(store *store.Store) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		rawSubmissionID := chi.URLParam(r, "submissionID")
		submissionID, err := strconv.ParseInt(rawSubmissionID, 10, 32)
		if err != nil {
			log.Printf("could not parse submission ID %s: %v\n", rawSubmissionID, err)
			writeBadRequest(w)
		}
		file, err := store.GetSubmissionFiles(int(submissionID))
		if err != nil {
			log.Printf("failed to get submissions file: %v\n", err)
			writeInternalServerError(w)
			return
		}
		w.Write(file)
	}))
}

func submissionFeedbackHandler(store *store.Store) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		// Unmarshal the POST body
		var feedbackReq SubmissionFeedbackRequest
		err := json.Unmarshal(request.GetBody(r), &feedbackReq)
		if err != nil {
			writeBadRequest(w)
			return
		}
		rawSubmissionID := chi.URLParam(r, "submissionID")
		submissionID, err := strconv.ParseInt(rawSubmissionID, 10, 32)
		if err != nil {
			log.Printf("could not parse submission ID %s: %v\n", rawSubmissionID, err)
			writeBadRequest(w)
		}
		resp, err := buildSubmissionsFeedbackResponse(store, int(submissionID), feedbackReq.Feedback)
		if err != nil {
			log.Printf("failed to build submissions feedback response: %v\n", err)
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}

func studentUploadHandler(store *store.Store) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got student upload!")
		rawCourseID, err := queryURL("course_id", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		courseID, err := strconv.Atoi(rawCourseID)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildStudentUploadResponse(store, courseID, request.GetBody(r))
		if err != nil {
			log.Println("Failed to build studentUploadResponse " + err.Error())
			writeInternalServerError(w)
			return
		}
		log.Println("Got upload respones: " + string(resp))
		w.Write(resp)
	}))
}

func submissionUploadHandler(store *store.Store, ci *ci.Ci) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		contentTypeParams, ok := r.Header["Content-Type"]
		if !ok {
			log.Println("submissionUploadHandler received bad content-type")
			writeBadRequest(w)
			return
		}
		if len(contentTypeParams) == 0 {
			log.Println("submissionUploadHandler received bad content-type")
			writeBadRequest(w)
			return
		}
		_, params, err := mime.ParseMediaType(contentTypeParams[0])
		if err != nil {
			log.Printf("failed to parse mediatype in submissionUploadHandler: %v\n", err)
			writeBadRequest(w)
			return
		}
		boundary, ok := params["boundary"]
		if !ok {
			log.Println("submissionUploadHandler could not access multipart-form boundary")
			writeBadRequest(w)
			return
		}
		body := request.GetBody(r)
		form, err := multipart.NewReader(bytes.NewReader(body), boundary).ReadForm(100 * mb)
		if err != nil {
			log.Printf("failed to read form file in submissionUploadHandler: %v\n", err)
			writeBadRequest(w)
			return
		}
		formValueParams, ok := form.Value["title"]
		if !ok || len(formValueParams) == 0 {
			log.Println("submissionUploadHandler could not find submission title")
			writeBadRequest(w)
		}
		title := formValueParams[0]
		formValueParams, ok = form.Value["description"]
		if !ok || len(formValueParams) == 0 {
			log.Println("submissionUploadHandler could not find submission description")
			writeBadRequest(w)
		}
		description := formValueParams[0]
		formValueParams, ok = form.Value["assessment_id"]
		if !ok || len(formValueParams) == 0 {
			log.Println("submissionUploadHandler could not find assessment_id")
			writeBadRequest(w)
		}
		rawAssessmentID := formValueParams[0]
		assessmentID, err := strconv.ParseInt(rawAssessmentID, 10, 64)
		if err != nil {
			log.Println("submissionUploadHandler could not parse assessment_id")
			writeBadRequest(w)
		}
		formValueParams, ok = form.Value["uid"]
		if !ok || len(formValueParams) == 0 {
			log.Println("submissionUploadHandler could not find uid")
			writeBadRequest(w)
		}
		uid := formValueParams[0]
		if _, ok := form.File["file"]; !ok {
			log.Println("submissionUploadHandler could not access uploaded file")
			writeBadRequest(w)
			return
		}
		if len(form.File["file"]) == 0 {
			log.Println("submissionUploadHandler could not access uploaded file")
			writeBadRequest(w)
			return
		}
		uploadedFile := form.File["file"][0]
		fileReader, err := uploadedFile.Open()
		if err != nil {
			log.Printf("failed to read file in submissionUploadHandler: %v\n", err)
			writeBadRequest(w)
			return
		}
		resp, err := buildSubmissionUploadResponse(store, ci, uid, title, description, int(assessmentID), fileReader)
		if err != nil {
			log.Println("Failed to build submissionUploadResponse " + err.Error())
			writeInternalServerError(w)
			return
		}
		log.Println("Got upload response: " + string(resp))
		w.Write(resp)
	}))
}

func logHandler(logger *log.Logger) http.HandlerFunc {
	return needsAuthorization(post(func(w http.ResponseWriter, r *http.Request) {
		// Read data somehow?
		// io.Copy(&buf, _)
		// log contents
		log.Println("Received remote logging")
	}))
}

func tutorialHandler(store store.Reader) http.HandlerFunc {
	return needsAuthorization(get(func(w http.ResponseWriter, r *http.Request) {
		uid, err := queryURL("uid", r)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := buildTutorialResponse(store, uid)
		if err != nil {
			log.Println("failed to build tutorial response")
			writeInternalServerError(w)
			return
		}
		w.Write(resp)
	}))
}
