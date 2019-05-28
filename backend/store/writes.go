package store

// WriteUser will upsert data in the database. If the user exists,
// (as defined by the uid), they will get new info. If they do not,
// a new user is created.
func (store *Store) WriteUser(user User) error {
	command := `
INSERT INTO users (uid, first_name, last_name, email, password) VALUES
	($1, $2, $3, $3, $4, $5)
ON CONFLICT DO UPDATE	
`
	_, err := store.db.Exec(command, user.UID, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	return err
}

// WriteTutorialEnrolment will enrol a user in a tutorial if they are not already enrolled.
func (store *Store) WriteTutorialEnrolment(uid string, tutorialID int) error {
	command := `
INSERT INTO tutorial_enrol (uid, tutorial_id) VALUES
	($1, $2)
ON CONFLICT DO UPDATE	
`
	_, err := store.db.Exec(command, uid, tutorialID)
	return err
}

// WriteCourseEnrolment will enrol a user in a course if they are not already enrolled
func (store *Store) WriteCourseEnrolment(uid string, role string, courseID int) error {
	command := `
INSERT INTO enrol (role, uid, course_id) VALUES
	($1, $2, $3)
ON CONFLICT DO UPDATE`
	_, err := store.db.Exec(command, role, uid, courseID)
	return err
}
