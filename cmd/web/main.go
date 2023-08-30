package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network addrress")
	flag.Parse()

	http.Handle("/", &home{})
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Print("Starting server on " + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
