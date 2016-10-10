package main

import (
	"errors"
	"github.com/twinj/uuid"
)

var (
	ErrUserNotFound          = errors.New("api: user not found")
	ErrUsernameAlreadyExists = errors.New("api: username already exists")
	ErrUsernameBlank         = errors.New("api: username cannot be empty")
	ErrPasswordBlank         = errors.New("api: username cannot be empty")
	ErrUserIdBlank           = errors.New("api: Id cannot be empty")
	ErrUserNil               = errors.New("api: user nil")
)

// User represents and user registered with the system.
type User struct {
	Id       string
	Username string
	Password string
}

// Create a new user object with a unique Id. Also checks for dupe username.
func NewUser(username, password string) (*User, error) {
	if username == "" {
		return nil, ErrUsernameBlank
	}
	if password == "" {
		return nil, ErrPasswordBlank
	}

	// Make sure there is no username conflict.
	_, err := db.GetUserByUsername(username)
	if err == ErrUserNotFound {
		return &User{Id: uuid.NewV1().String(), Username: username, Password: password}, nil
	} else if err != nil {
		return nil, err
	}
	// If no error, then a user already exists with that name.
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
