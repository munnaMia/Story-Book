// All helper methods
package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
