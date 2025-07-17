package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a http.Handler our a pointer to the servemux containing our application routes.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Handle 404 not found.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//create a fileserver. for serving static files as a http handler form the root of the application.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/blog/view/:id", app.blogView)
	router.HandlerFunc(http.MethodGet, "/blog/create", app.blogCreate)
	router.HandlerFunc(http.MethodPost, "/blog/create", app.blogCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
