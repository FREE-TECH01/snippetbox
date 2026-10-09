package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.ffestus.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	//command line flag for network  address
	addr := flag.String("addr", ":4000", "HTTP network address")

	//define a new command line for mysql DSN string
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=True", "mysql data source name")

	//parsing command line flag
	flag.Parse()

	//starting the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//opening the DB connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//closing connection before main function exists
	defer db.Close()

	//initializing the instance of our application struct, containing the dependencies
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	//using the info method to log the starting server at the info severity
	logger.Info("starting server", "addr", *addr)

	//using the Error method to log any error message returned by the server
	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
