package main

import (
	"fmt"
	"log"
	"net/http"
)

const version = "0.0.1"

func main() {
	fmt.Println("JAMA Server version ", version)
	fmt.Println("[*] Starting server...")

	StartServer()
}

func StartServer() {
	http.HandleFunc("/", HomeHandler)            // Return index.html
	http.HandleFunc("/favicon.ico", iconHandler) // Return favicon

	http.HandleFunc("/api/", apiHandler)               // Return API reference
	http.HandleFunc("/api/user", apiUserHandler)       // Handle all user API requests
	http.HandleFunc("/api/message", apiMessageHandler) // Handle all message API requests
	http.HandleFunc("/api/auth", apiAuthHandler)       // Handle all auth API requests

	log.Fatal(http.ListenAndServe(":8080", nil))
	return
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving index.html")
	filename := "index.html"
	http.ServeFile(w, r, filename)
}

func iconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving favicon.ico")
	filename := "favicon.ico"
	http.ServeFile(w, r, filename)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving api.html")
	filename := "api.html"
	http.ServeFile(w, r, filename)
}
