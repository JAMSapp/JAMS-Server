package main

import (
	"net/http"
)

func apiAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Allow", "POST")
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
