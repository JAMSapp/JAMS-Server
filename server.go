package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

var db *bolt.DB

func main() {
	fmt.Println("JAMA Server version ", version)
	fmt.Println("[*] Starting server...")

	StartServer()
}

func StartServer() {
	fmt.Println("[+] Loading BoltDB")
	d, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("[!] Error opening BoltDB database file")
		os.Exit(1)
	}
	db = d
	fmt.Println("[+] BoltDB loaded")

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
