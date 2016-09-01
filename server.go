package main

import (
	"log"
	"net/http"
)

const version = "0.0.1"

func main() {
	http.HandleFunc("/", HomeHandler) // Return about page

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	filename := "index.html"
	http.ServeFile(w, r, filename)
}
