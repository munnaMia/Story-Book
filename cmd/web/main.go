package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", blogView)
	mux.HandleFunc("/blog/create", blogCreate)

	log.Println("Server running at PORT: 8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
