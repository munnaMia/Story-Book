package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/munnaMia/Story-Book/internal/model"
	"github.com/munnaMia/Story-Book/internal/validator"
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
// The struct tag `form:"-"`
// tells the decoder to completely ignore a field during decoding.
type blogCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = blogCreateForm{
		Expires: 365,
	}

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
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	var form blogCreateForm

	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct. This will
	// essentially fill our struct with the relevant values from the HTML form.
	// If there is a problem, we return a 400 Bad Request response to the client.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// pass data to insert method
	id, err := app.blogs.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
