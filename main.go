package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startHandler)          // GET /start
	http.HandleFunc("/start-game", startPostHandler) // POST /start-game
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
