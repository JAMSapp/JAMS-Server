package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrUserNotFound         = errors.New("api: user not found")
	ErrUserAlreadyExists    = errors.New("api: user already exists")
	ErrUnsupportedMediaType = errors.New("api: Content-Type unsupported")
)

func apiUserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// TODO: Factor out get all users
	if id == "" {
		users, err := db.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		buf, err := json.Marshal(users)
		fmt.Fprintf(w, "%s", string(buf))
		return
	}
	user, err := db.GetUserById(id)
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the user as a JSON object and 200 OK.
	fmt.Fprintf(w, "%s", string(buf))
}

// Handle POST requests at the /api/user URL.
func apiUserPostHandler(w http.ResponseWriter, r *http.Request) {
	// Require JSON content type.
	content := r.Header.Get("Content-Type")
	if content != "application/json" {
		http.Error(w, ErrUnsupportedMediaType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a temp User to unmarshal
	var temp User
	err = json.Unmarshal(body, &temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new user with a randomly generated ID.
	user := NewUser(temp.Username, temp.Password)

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
	err = db.SaveUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal newly created user into JSON for response
	buf, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure we return 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", string(buf))
	return
}

// PUT /api/user/{id} - Update a user's record or create one if it doesn't exist
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

	// Check if a user exists.
	_, err = db.GetUserById(user.Id)
	if err != ErrUserNotFound {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	err = db.SaveUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func apiUserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	u, err := db.GetUserById(id)
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Possible for user to be deleted between these, though highly unlikely.
	err = u.Delete()
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	return
}
