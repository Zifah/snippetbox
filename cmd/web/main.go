package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr         string
	praiseAuthor bool
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	conf := config{}
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network addrress")
	flag.BoolVar(&conf.praiseAuthor, "praiseAuthor", false, "Praise or demean author")
	flag.Parse()

	app := application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	if conf.praiseAuthor {
		app.infoLog.Print("Hafiz is a great author!")
	} else {
		app.infoLog.Print("Hafiz is a mediocre author!")
	}

	srv := http.Server{
		Addr:     conf.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", conf.addr)
	app.errorLog.Fatal(srv.ListenAndServe())
}
