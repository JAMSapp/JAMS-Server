package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TODO: factor out magic numbers and strings

func TestRoutes(t *testing.T) {
	go StartServer()
	time.Sleep(100 * time.Millisecond)

	testMessage(t)
	testUser(t)
	testAuth(t)
}

func testMessage(t *testing.T) {
	r := Get("/api/message", t)
	if r != 200 {
		t.Errorf("GET /api/message did not return 200, instead returned %d\n", r)
	}
}

func testUser(t *testing.T) {
	r := Get("/api/user", t)
	if r != 200 {
		t.Errorf("GET /api/user did not return 200, instead returned %d\n", r)
	}

	r = Post("/api/user", "{\"id\": 123, \"username\":\"asdf\", \"password\": \"fdsa\"}", t)
	if r != 201 && r != 409 {
		t.Errorf("POST /api/user did not return 201, instead returned %d\n", r)
	}

	r = Get("/api/user/123", t)
	if r != 200 {
		t.Errorf("GET /api/user/123 did not return 200, instead returned %d\n", r)
	}

	r = Delete("/api/user/123", t)
	if r != 204 {
		if r == 404 {
			t.Errorf("DELETE on user not found\n")
		}
		t.Errorf("DELETE did not return 204, instead returned %d\n", r)
	}
}

func testAuth(t *testing.T) {
	r := Get("/api/auth", t)
	if r != 405 {
		t.Errorf("GET /api/auth did not return 405, instead returned %d\n", r)
	}
}

func Get(path string, t *testing.T) int {
	get_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	res, err := http.Get(get_url)
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}
	defer res.Body.Close()

	fmt.Printf("GET %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}

func Post(path, body string, t *testing.T) int {
	post_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	bodyType := "application/json"
	res, err := http.Post(post_url, bodyType, strings.NewReader(body))
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}
	defer res.Body.Close()

	fmt.Printf("POST %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}

func Put(path, body string, t *testing.T) int {
	put_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	bodyType := "application/json"
	req, err := http.NewRequest("PUT", put_url, strings.NewReader(body))
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}

	req.Header.Add("Content-Type", bodyType)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}
	defer res.Body.Close()

	fmt.Printf("PUT %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}

func Delete(path string, t *testing.T) int {
	del_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	req, err := http.NewRequest("DELETE", del_url, strings.NewReader(""))
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf(err.Error())
		return 0
	}
	defer res.Body.Close()

	fmt.Printf("DELETE %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}
