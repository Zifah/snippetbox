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

func main() {
	conf := config{}
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network addrress")
	flag.BoolVar(&conf.praiseAuthor, "praiseAuthor", false, "Praise or demean author")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	http.Handle("/", &home{})
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	if conf.praiseAuthor {
		log.Print("Hafiz is a great author!")
	} else {
		log.Print("Hafiz is a mediocre author!")
	}

	infoLog.Printf("Starting server on %s", conf.addr)
	errorLog.Fatal(http.ListenAndServe(conf.addr, nil))
}
