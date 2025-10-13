package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var game *Game

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if game == nil {
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

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

func startHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/start.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func startPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		player1 := r.FormValue("player1")
		player2 := r.FormValue("player2")
		level := r.FormValue("level")
		game = NewGame(level, player1, player2)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err == nil && game != nil {
			game.Play(col)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		game = nil
	}
	http.Redirect(w, r, "/start", http.StatusSeeOther)
}
