package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler) // <-- Place cette ligne ici

	http.ListenAndServe(":8080", nil)
}
