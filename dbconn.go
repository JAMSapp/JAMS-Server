package main

const (
	USERS    = "users"
	MESSAGES = "messages"
	UNREAD   = "unread_messages"
)

type DBConn interface {
	SaveUser(u *User) error
	DeleteUser(u *User) error
	GetUserById(id string) (*User, error)
	GetUsers() ([]User, error)
	SaveMessage(m *Message) error
	DeleteMessage(m *Message) error
	AddUnreadMessage(u *User, m *Message) error
}
