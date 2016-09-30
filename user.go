package main

import (
	"errors"
	"github.com/twinj/uuid"
)

var (
	ErrUserNotFound          = errors.New("api: user not found")
	ErrUsernameAlreadyExists = errors.New("api: username already exists")
	ErrUserObjectNil         = errors.New("api: user nil")
)

// User represents and user registered with the system.
type User struct {
	Id       string
	Username string
	Password string
}

func NewUser(username, password string) (*User, error) {
	_, err := db.GetUserByUsername(username)
	if err == ErrUserNotFound {
		return &User{Id: uuid.NewV1().String(), Username: username, Password: password}, nil
	} else if err != nil {
		return nil, err
	}
	return nil, ErrUsernameAlreadyExists
}

func (u *User) Save() error {
	existing, err := db.GetUserByUsername(u.Username)

	if err != nil {
		// If User doesn't yet exist, we're good to go
		if err == ErrUserNotFound {
			return db.SaveUser(u)
		}
		// Some other error is no bueno
		return err
	}

	// If the username belongs to a different User ID we fail
	if existing.Id != u.Id {
		return ErrUsernameAlreadyExists
	}

	// Otherwise they should match and we should be good.
	return db.SaveUser(u)
}

func (u *User) Delete() error {
	return db.DeleteUser(u)
}

func (u *User) SendMessage(m *Message) error {
	return db.AddUnreadMessage(u, m)
}
