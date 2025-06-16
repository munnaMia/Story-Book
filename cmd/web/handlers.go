package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	/*
		r.URL.Path != "/"
		=================

		"/" - this route work as a wildcard which give response to any route that user hit
		to prevent this happen we use this line.
	*/
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Storybook"))
}

func blogView(w http.ResponseWriter, r *http.Request) {
	// id which given by user should be a int and bigger then 0.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Displaying a specific blog by id : %d", id)
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		/*
			r.Method != "POST"
			==================

			This block of code check the method and return a 405 status code with a massage
			"Method Not Allowed" using w.WriteHeader(http.StatusMethodNotAllowed) & w.write("Method not allowed")

			Learn more about
				- https://pkg.go.dev/net/http#pkg-constants
				- https://pkg.go.dev/net/http#DetectContentType


			!) It’s only possible to call w.WriteHeader() once per response, and after the
			status code has been written it can’t be changed. If you try to call w.WriteHeader()
			a second time Go will log a warning message.

			1-- first sample:
			-----------------
				w.Header().Set("Allow", "POST")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method Not Allowed"))
		*/
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new blog..."))
}
