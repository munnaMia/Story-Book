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
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/munnaMia/Story-Book/internal/model"
)

// application struct to hold the application-wide dependencies for the web application.
// lower case struct name for internal use
type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	blogs    *model.BlogModel
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
		Note: A quirk of our MySQL driver is that we need to use the parseTime=true parameter
		in our DSN to force it to convert TIME and DATE fields to time.Time. Otherwise it returns
		these as []byte objects. This is one of the many driver-specific parameters that it offers.
	*/
	dsn := flag.String("dsn", "webhost:pass@/storybook?parseTime=true", "MySQL data source name")

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

	/*
		To keep the main() function tidy I've put the code for creating a connection
		pool into the separate openDB() function below. We pass openDB() the DSN
		from the command-line flag.
	*/
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new instance of our application struct, containing the dependencies.
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		blogs:    &model.BlogModel{DB: db},
	}

	/*
		set	the ErrorLog field so that the server now uses the custom errorLog logger in
		the event of any problems.
	*/
	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Previously we done this in this way : log.Printf("Server running at PORT: %s \n", *addr)
	infoLog.Printf("Server running at PORT: %s \n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
