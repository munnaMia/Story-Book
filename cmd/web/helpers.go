// All helper methods
package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	/*
		In the serverError() helper we use the debug.Stack() function to get a stack trace
		for the current goroutine and append it to the log message.
		Being able to see the execution path of the application via the stack trace can be helpful
		when you’re trying to debug errors.
	*/
	trace := fmt.Sprintf("%s \n %s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // depth 2 reason to find where the error occur from the stack.

	// e http.StatusText() function to automatically generate a human-friendly text representation of a given HTTP status code
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	/*
		Retrieve the appropriate template set from the cache based on the page
		name (like 'home.tmpl'). If no entry exists in the cache with the
		provided name, then create a new error and call the serverError() helper
		method that we made earlier and return.
	*/
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
	}

	// Initialize the buffer
	buf := new(bytes.Buffer)

	/*
		Write the template to the buffer, instead of straight to the
		http.ResponseWriter. If there's an error, call our serverError() helper
		and then return.
	*/
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	/*
		Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
		If the template is written to the buffer without any errors, we are safe
		to go ahead and write the HTTP status code to http.ResponseWriter.
	*/
	w.WriteHeader(status)

	/*
		Write the contents of the buffer to the http.ResponseWriter. Note: this
		is another time where we pass our http.ResponseWriter to a function that
		takes an io.Writer.*
	*/
	buf.WriteTo(w)

}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year. Note that we're not using the
// *http.Request parameter here at the moment, but we will do later.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return err
}
