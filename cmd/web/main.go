package main

import (
	"flag"
	"log"
	"net/http"
)

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

	mux := http.NewServeMux()

	//create a fileserver. for serving static files as a http handler form the root of the application.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", blogView)
	mux.HandleFunc("/blog/create", blogCreate)

	log.Printf("Server running at PORT: %s \n", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
