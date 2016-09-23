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

	// Test message api
	testMessage(t)

	// Test user api
	testUser(t)

	// Test auth api
	testAuth(t)
}

func testMessage(t *testing.T) {
	// GET /api/message
	// TODO: Decide on proper response and functionality for GET /api/message
	r := Get("/api/message", t)
	if r != 200 {
		t.Errorf("GET /api/message did not return 200, instead returned %d\n", r)
	}

	// POST /api/message
	// Creates a new message with the given body parameter.
	// Should return 201 on successful creation or 5** for server error.
	r = Post("/api/message", "{\"body\": \"This is the body of the first message.\"}", t)
	if r != 201 {
		t.Errorf("POST /api/message did not return 201, instead returned %d\n", r)
	}
}

func testUser(t *testing.T) {
	// GET /api/user
	// TODO: Decide on use of this call. Is it needed?
	r := Get("/api/user", t)
	if r != 200 {
		t.Errorf("GET /api/user did not return 200, instead returned %d\n", r)
	}

	// PUT /api/user/{id}
	// Create a new user or update a user record.
	// Should return 201 on successful creation, 200 on update, or 5** for server error.
	r = Put("/api/user/321", "{\"id\": 321, \"username\":\"fdsa\", \"password\": \"fdsa\"}", t)
	if r != 201 {
		t.Errorf("PUT /api/user/321 did not return 201, instead returned %d\n", r)
	}

	// GET /api/user/{id}
	// Returns the given user resource
	// Should return 200 if found or 404 if resource isn't found.
	r = Get("/api/user/321", t)
	if r != 200 {
		t.Errorf("GET /api/user/321 did not return 200, instead returned %d\n", r)
	}

	// POST /api/user
	// Create a new user
	// TODO: Remove id from this and generate a unique one automatically.
	// Should return 201 for created or 409 for conflict.
	r = Post("/api/user", "{\"id\": 123, \"username\":\"asdf\", \"password\": \"fdsa\"}", t)
	if r != 201 && r != 409 {
		t.Errorf("POST /api/user did not return 201, instead returned %d\n", r)
	}

	// Get newly created user.
	r = Get("/api/user/123", t)
	if r != 200 {
		t.Errorf("GET /api/user/123 did not return 200, instead returned %d\n", r)
	}

	// Update the user's record.
	r = Put("/api/user/123", "{\"id\": 123, \"username\":\"fdsa\", \"password\": \"fdsa\"}", t)
	if r != 200 {
		t.Errorf("PUT /api/user/123 did not return 200, instead returned %d\n", r)
	}

	// DELETE /api/user/{id}
	// Delete's the user resource.
	// Should return 204 for deletion, 404 if not found, or 5** for server error
	r = Delete("/api/user/123", t)
	if r != 204 {
		t.Errorf("DELETE did not return 204, instead returned %d\n", r)
	}

	// Delete user record.
	r = Delete("/api/user/321", t)
	if r != 204 {
		t.Errorf("DELETE did not return 204, instead returned %d\n", r)
	}

	// Test delete on non existant user
	r = Delete("/api/user/1337", t)
	if r != 404 {
		t.Errorf("DELETE on user not found\n")
	}
}

func testAuth(t *testing.T) {
	// GET /api/auth
	// TODO: Will return an auth token or something eventually...
	// Should return 405 unauthorized until implemented.
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
