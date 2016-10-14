package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func apiGetAllThreadsHandler(w http.ResponseWriter, r *http.Request) {
	threads, err := db.GetAllThreads()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(threads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the user as a JSON object and 200 OK.
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(buf))
}

func apiThreadGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	thread, err := db.GetThread(id)
	if err != nil {
		if err == ErrThreadNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user as a JSON object and 200 OK.
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(buf))
}

func apiThreadMessageGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	thread, err := db.GetThread(id)
	if err != nil {
		if err == ErrThreadNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := db.GetThreadMessages(thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user as a JSON object and 200 OK.
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(buf))
}

func apiThreadMessagePostHandler(w http.ResponseWriter, r *http.Request) {
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
	var temp Message
	err = json.Unmarshal(body, &temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	thread_id := vars["id"]

	thread, err := db.GetThread(thread_id)
	if err != nil {
		if err == ErrThreadNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mes := NewMessage(temp.Body)
	err = db.SaveMessage(mes, thread)
	if err != nil {
		if err == ErrMsgBodyBlank {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal newly created user into JSON for response
	buf, err := json.Marshal(mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure we return 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", string(buf))
	return
}
