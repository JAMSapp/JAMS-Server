package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrMessageNotFound = errors.New("api: message not found")
)

func apiThreadMessagePostHandler(w http.ResponseWriter, r *http.Request) {
	content := r.Header.Get("Content-Type")
	if content != "application/json" {
		http.Error(w, ErrUnsupportedMediaType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mes := NewMessage("")
	err = json.Unmarshal(body, &mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = mes.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
func apiMessageGetHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := db.GetMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(buf))
}

// For anything not a GET.
func apiMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Allowed", "GET")
	http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	return
}
