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
	r.HandleFunc("/api/user/{id}", apiUserGetHandler).Methods("GET")                // Get a user
	r.HandleFunc("/api/user/{id}", apiUserPutHandler).Methods("PUT")                // Update a user
	r.HandleFunc("/api/user/{id}", apiUserDeleteHandler).Methods("DELETE")          // Delete a user
	r.HandleFunc("/api/user", apiUserPostHandler).Methods("POST")                   // Create a new user
	r.HandleFunc("/api/user", apiUserGetHandler).Methods("GET")                     // Get all users
	r.HandleFunc("/api/user/{id}/thread", apiUserThreadGetHandler).Methods("GET")   // Get a user's threads
	r.HandleFunc("/api/user/{id}/thread", apiUserThreadPostHandler).Methods("POST") // Create a new message thread

	r.HandleFunc("/api/thread", apiGetAllThreadsHandler).Methods("GET")                   // Get all threads
	r.HandleFunc("/api/thread/{id}", apiThreadGetHandler).Methods("GET")                  // Get a thread's info
	r.HandleFunc("/api/thread/{id}/message", apiThreadMessageGetHandler).Methods("GET")   // Get all messages for this thread
	r.HandleFunc("/api/thread/{id}/message", apiThreadMessagePostHandler).Methods("POST") // Create a new message for this thread

	// /message
	r.HandleFunc("/api/message", apiMessageGetHandler).Methods("GET") // Get all messages.
	r.HandleFunc("/api/message", apiMessageHandler)                   // Method not allowed

	// /auth
	r.HandleFunc("/api/auth", apiAuthHandler) // Handle all auth API requests

	return r
}
