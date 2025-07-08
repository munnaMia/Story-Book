package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/munnaMia/Story-Book/internal/model"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		r.URL.Path != "/"
		=================

		"/" - this route work as a wildcard which give response to any route that user hit
		to prevent this happen we use this line.
	*/
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	blogs, err := app.blogs.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := templateData{
		Blogs: blogs,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	// id which given by user should be a int and bigger then 0.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/view.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// instance of a templateData struct holding the snippet data.
	data := &templateData{
		Blog: blog,
	}

	// HTML templates, any dynamic data that you pass in is represented by the .
	// character (referred to as dot).
	// In this specific case, the underlying type of dot will be a models.Blog struct. When the
	// underlying type of dot is a struct, you can render (or yield) the value of any exported field in
	// your templates by postfixing dot with the field name. So, because our models.Snippet struct
	// has a Title field, we could yield the snippet title by writing {{.Title}} in our templates.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
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
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//DUMY DATA for test purpose.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// pass dumy data to insert method to test
	id, err := app.blogs.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/blog/view?id=%d", id), http.StatusSeeOther)
}
