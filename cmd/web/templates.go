package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/munnaMia/Story-Book/internal/model"
)

/*
Define a templateData type to act as the holding structure for
any dynamic data that we want to pass to our HTML templates.
At the moment it only contains one field, but we'll add more
to it as the build progresses.
*/
type templateData struct {
	CurrentYear int
	Blog        *model.Blog
	Blogs       []*model.Blog
	Form        any
	Flash       string
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	/*
		Use the filepath.Glob() function to get a slice of all filepaths that
		match the pattern "./ui/html/pages/*.html". This will essentially gives
		us a slice of all the filepaths for our application 'page' templates
		like: [ui/html/pages/home.html ui/html/pages/view.html]
	*/
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.html') as the key.
		cache[name] = ts
	}

	// return the map
	return cache, nil
}
