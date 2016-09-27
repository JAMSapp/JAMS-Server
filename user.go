package main

import (
	"github.com/twinj/uuid"
)

// User represents and user registered with the system.
type User struct {
	Id       string
	Username string
	Password string
}

func NewUser(username, password string) *User {
	return &User{Id: uuid.NewV1().String(), Username: username, Password: password}
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
