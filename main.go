package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w,r)
		return
	}
	w.Write([]byte("Welcome to Storybook"))
}

func blogView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("View a blog"))
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new blog..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", blogView)
	mux.HandleFunc("/blog/create", blogCreate)

	log.Println("Server running at PORT: 8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
