package main

import (
	"fmt"
	"github.com/braintree/manners"
	"net/http"
	"os"
)

const version = "0.0.1"

var db DBConn

func main() {
	fmt.Println("JAMA Server version ", version)
	fmt.Println("[*] Starting server...")

	c := make(chan int)
	go StartServer(c)
	<-c
}

func StartServer(c chan int) {
	fmt.Println("[+] Loading BoltDB")
	boltdb, err := BoltDBOpen(DBFILE)
	if err != nil {
		fmt.Printf("[!] Error opening BoltDB: %s\n", err.Error())
		os.Exit(1)
	}
	defer boltdb.Close()
	db = boltdb
	fmt.Println("[+] BoltDB loaded")

	r := routes()

	manners.ListenAndServe(":8080", r)
	c <- 0
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
