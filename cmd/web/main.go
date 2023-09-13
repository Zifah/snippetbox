package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr         string
	praiseAuthor bool
	dsn          string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	conf := config{}
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network addrress")
	flag.BoolVar(&conf.praiseAuthor, "praiseAuthor", false, "Praise or demean author")
	flag.StringVar(&conf.dsn, "dsn", "web:Password123@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(conf.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	if conf.praiseAuthor {
		infoLog.Print("Hafiz is a great author!")
	} else {
		infoLog.Print("Hafiz is a mediocre author!")
	}

	srv := http.Server{
		Addr:     conf.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", conf.addr)
	errorLog.Fatal(srv.ListenAndServe())
}
