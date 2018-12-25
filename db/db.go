package db

// User represents a user in the database.
// It can be a student, tutor, convenor or
// admin.
type User struct {
	Name  string
	UID   string
	Roles []string
}

// Student represents a student in the database.
// It is guaranteed to contain a "student" role.
type Student struct {
	User
}

// Reader is the interface for read-only interactions
// with the database. It is implemented by Store.
type Reader interface {
	GetUser(uid string) (User, error)
	//GetStudent(uid string) (Student, error)
}

// Writer is the interface for write-only interactions
// with the database. It is implemented by Store.
type Writer interface {
	//CreateUser(user User) error
	//SetUser(user User) error
}
