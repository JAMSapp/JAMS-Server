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
	http.HandleFunc("/", HomeHandler)            // Return about page
	http.HandleFunc("/favicon.ico", iconHandler) // Return about page

	http.HandleFunc("/api/user", apiHandler)    // Return about page
	http.HandleFunc("/api/message", apiHandler) // Return about page
	http.HandleFunc("/api/auth", apiHandler)    // Return about page

	log.Fatal(http.ListenAndServe(":8080", nil))
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
