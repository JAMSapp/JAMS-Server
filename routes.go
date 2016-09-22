package main

import (
	"fmt"

	"github.com/gorilla/mux"
)

func routes() *mux.Router {
	fmt.Println("[+] Setting Routes...")

	r := mux.NewRouter()

	// Static content the lives in public folder
	r.HandleFunc("/", publicHandler)
	r.HandleFunc("/favicon.ico", publicHandler)

	// Base API functions
	r.HandleFunc("/api/", apiHandler) // Return API reference

	// /user
	r.HandleFunc("/api/user/{id}", apiUserGetHandler).Methods("GET")       // Get a user
	r.HandleFunc("/api/user/{id}", apiUserPutHandler).Methods("PUT")       // Update a user
	r.HandleFunc("/api/user/{id}", apiUserDeleteHandler).Methods("DELETE") // Delete a user
	r.HandleFunc("/api/user", apiUserPostHandler).Methods("POST")          // Create a new user
	r.HandleFunc("/api/user", apiUserGetHandler).Methods("GET")            // Get all users

	// /message
	r.HandleFunc("/api/message", apiMessageHandler) // Handle all message API requests

	// /auth
	r.HandleFunc("/api/auth", apiAuthHandler) // Handle all auth API requests

	return r
}
