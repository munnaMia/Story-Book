package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/munnaMia/Story-Book/internal/model"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		r.URL.Path != "/"
		=================

		"/" - this route work as a wildcard which give response to any route that user hit
		to prevent this happen we use this line.

		update with httprouter pkg:
		if r.URL.Path != "/" {
			app.notFound(w)
			return
		}
	*/

	blogs, err := app.blogs.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	param := httprouter.ParamsFromContext(r.Context())

	// id which given by user should be a int and bigger then 0.
	id, err := strconv.Atoi(param.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// blog data will be render on html
	blog, err := app.blogs.Get(id)

	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Blog = blog

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) blogCreatePost(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	/*
	// 		r.Method != "POST"
	// 		==================

	// 		This block of code check the method and return a 405 status code with a massage
	// 		"Method Not Allowed" using w.WriteHeader(http.StatusMethodNotAllowed) & w.write("Method not allowed")

	// 		Learn more about
	// 			- https://pkg.go.dev/net/http#pkg-constants
	// 			- https://pkg.go.dev/net/http#DetectContentType

	// 		!) It’s only possible to call w.WriteHeader() once per response, and after the
	// 		status code has been written it can’t be changed. If you try to call w.WriteHeader()
	// 		a second time Go will log a warning message.

	// 		1-- first sample:
	// 		-----------------
	// 			w.Header().Set("Allow", "POST")
	// 			w.WriteHeader(http.StatusMethodNotAllowed)
	// 			w.Write([]byte("Method Not Allowed"))
	// 	*/
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//DUMY DATA for test purpose.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// pass dumy data to insert method to test
	id, err := app.blogs.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
