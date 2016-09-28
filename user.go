package main

import (
	"errors"
	"github.com/twinj/uuid"
)

var (
	ErrUserNotFound      = errors.New("api: user not found")
	ErrUserAlreadyExists = errors.New("api: user already exists")
	ErrUserObjectNil     = errors.New("api: user nil")
)

// User represents and user registered with the system.
type User struct {
	Id       string
	Username string
	Password string
}

func NewUser(username, password string) (*User, error) {
	_, err := db.GetUserByUsername(username)
	if err != ErrUserNotFound {
		return nil, ErrUserAlreadyExists
	}
	return &User{Id: uuid.NewV1().String(), Username: username, Password: password}, nil
}

func (u *User) Save() error {
	return db.SaveUser(u)
}

func (u *User) Delete() error {
	return db.DeleteUser(u)
}

func (u *User) SendMessage(m *Message) error {
	return db.AddUnreadMessage(u, m)
}
