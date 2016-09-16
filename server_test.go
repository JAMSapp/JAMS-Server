package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go StartServer()
	time.Sleep(100 * time.Millisecond)
	r := Get("/api/auth", t)
	if r != 405 {
		t.Errorf("GET /api/auth did not return 405, instead returned %d\n", r)
	}

	r = Get("/api/user", t)
	if r != 404 {
		t.Errorf("GET /api/user did not return 200, instead returned %d\n", r)
	}

	r = Get("/api/message", t)
	if r != 200 {
		t.Errorf("GET /api/message did not return 200, instead returned %d\n", r)
	}

}

func Get(path string, t *testing.T) int {
	get_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	res, err := http.Get(get_url)
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}

	fmt.Printf("GET %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}