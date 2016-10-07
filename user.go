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

// Create a new user object with a unique Id. Also checks for dupe username.
func NewUser(username, password string) (*User, error) {
	_, err := db.GetUserByUsername(username)
	if err == ErrUserNotFound {
		return &User{Id: uuid.NewV1().String(), Username: username, Password: password}, nil
	} else if err != nil {
		return nil, err
	}
	return nil, ErrUsernameAlreadyExists
}

// Save a User object to the DB. Enforces username uniqueness
func (u *User) Save() error {
	return db.SaveUser(u)
}

// Delete a user from the database.
func (u *User) Delete() error {
	return db.DeleteUser(u)
}
