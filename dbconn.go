package main

const (
	USERS    = "users"
	MESSAGES = "messages"
)

type DBConn interface {
	SaveUser(u *User) error
	DeleteUser(u *User) error
	GetUserById(id int) (*User, error)
	GetUsers() ([]User, error)
	SaveMessage(m *Message) error
	DeleteMessage(m *Message) error
}
