package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
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
	// Checks for username conflict as well.
	user, err := NewUser(temp.Username, temp.Password)
	if err == ErrUsernameAlreadyExists {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Must save for persistence.
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

	// Get the user data from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON into a user
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if a user exists with that ID (update or creation?)
	var newUser bool
	_, err = db.GetUserById(user.Id)

	if err != nil {
		// If there was no user by this ID we return 201
		if err == ErrUserNotFound {
			newUser = true

		} else { // If we got a different error for some reason, return 5**
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else { // No error, there is an existing user so this is just an update
		newUser = false
	}

	err = user.Save()
	// Check to see if there was an error saving this user.
	if err != nil {
		// If the username already belongs to another User ID we should return
		// conflict
		if err == ErrUsernameAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		// If it was a different error, this is a server fault.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if newUser {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
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
