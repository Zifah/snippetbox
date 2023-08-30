package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", &home{})
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Print("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
