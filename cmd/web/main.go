package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network addrress")
	printAuthor := flag.Bool("printAuthor", false, "Print author name")
	flag.Parse()

	http.Handle("/", &home{})
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	if *printAuthor {
		log.Print("Hafiz is a great author!")
	} else {
		log.Print("Hafiz is a mediocre author!")
	}

	log.Print("Starting server on " + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
