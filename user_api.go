package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	ErrUserNotFound         = errors.New("api: user not found")
	ErrUserAlreadyExists    = errors.New("api: user already exists")
	ErrUnsupportedMediaType = errors.New("api: Content-Type unsupported")
)

func apiUserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid := vars["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		// If it doesn't parse to an int, there will be no associated user.
		http.Error(w, err.Error(), http.StatusNotFound)
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

func apiUserPutHandler(w http.ResponseWriter, r *http.Request) {
	content := r.Header.Get("Content-Type")
	if content != "application/json" {
		http.Error(w, ErrUnsupportedMediaType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	var user User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Double check that a user doesn't already exist.
	_, err = db.GetUserById(user.Id)
	if err != ErrUserNotFound {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, ErrUserAlreadyExists.Error(), http.StatusConflict)
		return
	}

	// TODO: Make a new call for saving a New user to avoid the above check.
	err = db.SaveUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
