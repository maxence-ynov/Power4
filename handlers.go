package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var game = NewGame()

// GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").
		Funcs(template.FuncMap{"seq": seq}).
		ParseFiles("templates/index.html"))

	tmpl.Execute(w, game)
}

// POST /play
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err == nil {
			game.Play(col)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Generates a sequence of integers from start to end inclusive
func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

// ...existing code...

// POST /reset
func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		game = NewGame()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ...existing code...
