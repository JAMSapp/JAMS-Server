package main

const (
	USERS    = "users"
	MESSAGES = "messages"
)

type DBConn interface {
	SaveUser(u *User) error
	//	GetUserById(id int) (*User, error)
}
