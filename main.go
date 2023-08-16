package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox."))
	// TODO Hafiz: Why is this a byte array and not a string?
	// Because the Write method only accepts a byte array
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet."))

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("View an existing snippet."))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
