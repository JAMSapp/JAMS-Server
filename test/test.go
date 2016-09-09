package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := Get("/api/auth")
	if r != 405 {
		fmt.Errorf("GET /api/auth did not return 405, instead returned %d\n", r)
	}

	r = Get("/api/user")
	if r != 200 {
		fmt.Errorf("GET /api/user did not return 200, instead returned %d\n", r)
	}

	r = Get("/api/message")
	if r != 200 {
		fmt.Errorf("GET /api/message did not return 200, instead returned %d\n", r)
	}

}

func Get(path string) int {
	get_url := fmt.Sprintf("http://127.0.0.1:8080%s", path)
	res, err := http.Get(get_url)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Printf("GET %s: %d\n", path, res.StatusCode)
	return res.StatusCode
}
