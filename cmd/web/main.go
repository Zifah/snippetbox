package main

import (
	"flag"
	"log"
	"net/http"
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

	log.Print("Starting server on " + conf.addr)
	log.Fatal(http.ListenAndServe(conf.addr, nil))
}
