package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var game = NewGame()

// GET /
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").
		Funcs(template.FuncMap{
			"seq": func(start, end int) []int {
				s := make([]int, end-start+1)
				for i := range s {
					s[i] = start + i
				}
				return s
			},
			"eq":  func(a, b int) bool { return a == b },
			"sub": func(a, b int) int { return a - b },
		}).
		ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /play
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err == nil {
			game.Play(col)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /reset
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		game = NewGame()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
