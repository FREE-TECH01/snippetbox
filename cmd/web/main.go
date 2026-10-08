package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	//command line flag for network  address
	addr := flag.String("addr", ":4000", "HTTP network address")

	//parsing command line flag
	flag.Parse()

	//starting the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//initializing the instance of our application struct, containing the dependencies
	app := &application{
		logger: logger,
	}

	//using the info method to log the starting server at the info severity
	logger.Info("starting server", "addr", *addr)

	//using the Error method to log any error message returned by the server
	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
