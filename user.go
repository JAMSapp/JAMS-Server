package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	ErrUserNotFound = errors.New("user: user not found")
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

func apiUserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid := vars["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := db.GetUserById(id)
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<h1>Id: %d</h1><div>Username: %s</br>Password: %s</div>", u.Id, u.Username, u.Password)
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
