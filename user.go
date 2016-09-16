package main

import (
	"errors"
	"net/http"
)

type User struct {
	Id       int
	Username string
	Password string
}

func (user *User) Save(db DBConn) error {
	if db == nil {
		return errors.New("user: DBConn is nil")
	}
	return db.SaveUser(user)
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
