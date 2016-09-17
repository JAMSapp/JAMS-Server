package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "0.0.1"

var db DBConn

func main() {
	fmt.Println("JAMA Server version ", version)
	fmt.Println("[*] Starting server...")

	StartServer()
}

func StartServer() {
	fmt.Println("[+] Loading BoltDB")
	boltdb, err := BoltDBOpen("my.db")
	if err != nil {
		fmt.Printf("[!] Error opening BoltDB: %s\n", err.Error())
		os.Exit(1)
	}
	db = boltdb
	fmt.Println("[+] BoltDB loaded")

	r := routes()
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
	return
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	filename := "public/" + r.URL.Path
	http.ServeFile(w, r, filename)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving api.html")
	filename := "public/api.html"
	http.ServeFile(w, r, filename)
}
