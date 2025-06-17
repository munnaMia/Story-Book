package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//create a fileserver. for serving static files as a http handler form the root of the application.
	fileserver:= http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", blogView)
	mux.HandleFunc("/blog/create", blogCreate)

	log.Println("Server running at PORT: 8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
