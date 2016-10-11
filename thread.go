package main

import (
	"errors"
	"github.com/twinj/uuid"
)

var (
	ErrThreadNil      = errors.New("api: thread nil")
	ErrThreadNotFound = errors.New("api: thread not found")
	ErrThreadIdBlank  = errors.New("api: thread id blank")
)

// Thread will provide a way to relate messages with users associated with a
// message thread. A message thread may have many participants, and a message
// thread will have several messages.
type Thread struct {
	Id      string
	UserIds []string
}

// NewThread takes a slice of strings, each supposedly a user id.
func NewThread(users []string) *Thread {
	return &Thread{Id: uuid.NewV1().String(), UserIds: users}
}

// Save will store the thread in the database and will return an error in case
// of DB failure.
func (t *Thread) Save() error {
	return db.SaveThread(t)
}

// Delete will remove any matching record from the database permanently and will
// return an error in case of DB failure.
func (t *Thread) Delete() error {
	return db.DeleteThread(t)
}
