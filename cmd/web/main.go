/*
	Simple Rule : try to avoid use fatal & panic outside of main() function.— it’s good practice to return errors instead, and only panic or exit
	directly from main().

	As I said above, my general recommendation is to log your output to standard streams and
	redirect the output to a file at runtime. But if you don’t want to do this, you can always open a
	file in Go and use it as your log destination. As a rough example:
	------------------------------------------------------------------------
	|	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	|	if err != nil {
	|		log.Fatal(err)
	|	}
	|	defer f.Close()
	|	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	------------------------------------------------------------------------

*/

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application struct to hold the application-wide dependencies for the web application. 
// lower case struct name for internal use 
type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
}

func main() {
	/*
		Importent Notes
		---------------
		Our application have some hardcoded data like addr:8080 and "static" routes.
		having this hardcoded value isn't ideal.There’s no separation between our
		configuration settings and code, and we can’t change the settings at runtime
		(which is important if you need different settings for development, testing and
		production environments).
	*/

	/*
		addr
		----
		it's a new command line flag. with a default value localhost:8080. and some help
		massasge to understand what this flag does.

		How to use it?
			EX --> go run .\cmd\web\. -addr=":8000"
	*/
	addr := flag.String("addr", "localhost:8080", "HTTP network address")

	/*
		Parse()
		-------
		this method will parse the command line flag. this read command line value and
		assing it to the addr / or else variable. this have to call before use addr in code base.
		otherwise it will use default 8080 value. if any error occur doing run time application
		will be terminated

		Note :
			use -help flag to get all flags info.

		Research :
			flag.StringVar and how to use it with struct to store cofigs data.
	*/
	flag.Parse()

	/*
		log.new()
		---------

		this is used to create custom logger. it take destination where to write , massage, other information like data and time
	*/
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of our application struct, containing the dependencies.
	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
	}

	mux := http.NewServeMux()

	//create a fileserver. for serving static files as a http handler form the root of the application.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/blog/view", app.blogView)
	mux.HandleFunc("/blog/create", app.blogCreate)

	/*
		set	the ErrorLog field so that the server now uses the custom errorLog logger in
		the event of any problems.
	*/
	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Previously we done this in this way : log.Printf("Server running at PORT: %s \n", *addr)
	infoLog.Printf("Server running at PORT: %s \n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
