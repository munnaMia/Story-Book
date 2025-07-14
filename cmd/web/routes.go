package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a http.Handler our a pointer to the servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//create a fileserver. for serving static files as a http handler form the root of the application.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/blog/view", app.blogView)
	mux.HandleFunc("/blog/create", app.blogCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
