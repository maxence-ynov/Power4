package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/play", PlayHandler)
	http.HandleFunc("/reset", ResetHandler)
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
