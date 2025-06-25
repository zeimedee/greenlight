package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/zeimedee/greenlight/internal/data"

	"github.com/zeimedee/greenlight/internal/jsonlog"
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
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// Define  an application struck to hold the dependencies for our http handlers, helpers and middleware.
// At the moment this only contains a copy of the config struct and a logger, but it will
// grow to include a lot mote as out build progresses.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	//Define an instance of the config struct
	var cfg config

	//Read value of the port and env command-line flags into the config struct. we
	//default to using port number 4000 and environment "development" if no corresponding flags are provided
	flag.IntVar(&cfg.port, "port", 4000, "API SERVER PORT")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost/greenlight?sslmode=disable", "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSql max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSql max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSql max connections connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per seconde")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	//Initialize a new logger  which writes messages to the standard out stream,
	//prefixed with the current date and time.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	//Declare an instance of the application struct, containing the config struct and the logger

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	//Declare a http server with some sensible timeout settings, which listens on the port provided in the config struct and uses the httprouter instance returned by app.routes() as the server handler

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//Start the http server.
	logger.PrintInfo("starting %s server on %s", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err = srv.ListenAndServe()

	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
