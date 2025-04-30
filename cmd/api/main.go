package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"log"
)

// Declare a string containing the application version number. later in the book
// we'll generate this automatically at build time, but got now we'll just store the version number
// as a hard coded global constant
const version = "1.0.0"

// define the config struct to hold all the configuration settings for out application.
// for now the only configuration setting will be the network port that we want the server to listen on,
// and the name of the current operating environment for the application(development, staging,production, etc) we
// will read the configuration settings from the command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// Define  an application struck to hold the dependencies for our http handlers, helpers and middleware.
// At the moment this only contains a copy of the config struct and a logger, but it will
// grow to include a lot mote as out build progresses.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	//Define an instance of the config struct
	var cfg config

	//Read value of the port and env command-line flags into the config struct. we
	//default to using port number 4000 and environment "development" if no corresponding flags are provided
	flag.IntVar(&cfg.port, "port", 4000, "API SERVER PORT")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost/greenlight?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	//Initialize a new logger  which writes messages to the standard out stream,
	//prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	//Declare an instance of the application struct, containing the config struct and the logger

	app := &application{
		config: cfg,
		logger: logger,
	}

	//Declare a http server with some sensible timeout settings, which listens on the port provided in the config struct and uses the httprouter instance returned by app.routes() as the server handler

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//Start the http server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
