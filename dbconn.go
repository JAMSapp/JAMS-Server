package main

const (
	USERIDS   = "userids"
	USERNAMES = "usernames"
	MESSAGES  = "messages"
	UNREAD    = "unread_messages"
)

type DBConn interface {
	Init() error
	Close()
	SaveUser(u *User) error
	DeleteUser(u *User) error
	GetUserById(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUsers() ([]User, error)
	SaveMessage(m *Message) error
	DeleteMessage(m *Message) error
	AddUnreadMessage(u *User, m *Message) error
	GetUnreadMessages(u *User) ([]Message, error)
}
