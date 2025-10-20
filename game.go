package main

type Game struct {
	Board          [][]int
	Rows           int
	Columns        int
	Turn           int
	Winner         int
	Draw           bool
	GravityInverse bool
	TurnCount      int
	Player1Name    string
	Player2Name    string
	Level          string
}

func NewGame(level string, player1Name, player2Name string) *Game {
	var rows, cols int
	switch level {
	case "easy":
		rows, cols = 6, 7
	case "normal":
		rows, cols = 6, 9
	case "hard":
		rows, cols = 7, 8
	default:
		rows, cols = 6, 7
	}

	board := make([][]int, cols)
	for i := range board {
		board[i] = make([]int, rows)
	}

	return &Game{
		Board:       board,
		Rows:        rows,
		Columns:     cols,
		Turn:        1,
		Winner:      0,
		Draw:        false,
		TurnCount:   0,
		Player1Name: player1Name,
		Player2Name: player2Name,
		Level:       level,
	}
}

func (g *Game) Play(col int) {
	if col < 0 || col >= g.Columns || g.Winner != 0 || g.Draw {
		return
	}

	row := -1
	if g.GravityInverse {
		for r := 0; r < g.Rows; r++ {
			if g.Board[col][r] == 0 {
				row = r
				break
			}
		}
	} else {
		for r := g.Rows - 1; r >= 0; r-- {
			if g.Board[col][r] == 0 {
				row = r
				break
			}
		}
	}

	if row == -1 {
		return
	}

	g.Board[col][row] = g.Turn
	g.TurnCount++

	row = g.applyGravity(col, row, g.GravityInverse)

	var gravityToggled bool
	if g.TurnCount%5 == 0 {
		g.GravityInverse = !g.GravityInverse
		row = g.applyGravity(col, row, g.GravityInverse)
		gravityToggled = true
	}

	if g.checkWin(col, row) {
		g.Winner = g.Turn
		return
	}

	if gravityToggled {
		if winner := g.scanForWinner(); winner != 0 {
			g.Winner = winner
			return
		}
	}

	if g.checkDraw() {
		g.Draw = true
		return
	}

	g.Turn = 3 - g.Turn
}

func (g *Game) applyGravity(lastCol, lastRow int, gravityInverse bool) int {
	type cell struct {
		value int
		row   int
	}

	for c := 0; c < g.Columns; c++ {
		stack := make([]cell, 0, g.Rows)
		for r := 0; r < g.Rows; r++ {
			if g.Board[c][r] != 0 {
				stack = append(stack, cell{value: g.Board[c][r], row: r})
				g.Board[c][r] = 0
			}
		}

		for i, entry := range stack {
			var newRow int
			if gravityInverse {
				newRow = i
			} else {
				newRow = g.Rows - 1 - i
			}
			g.Board[c][newRow] = entry.value
			if c == lastCol && entry.row == lastRow {
				lastRow = newRow
			}
		}
	}

	return lastRow
}

func (g *Game) checkWin(col, row int) bool {
	player := g.Board[col][row]
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonale descendante
		{1, -1}, // Diagonale montante
	}

	for _, d := range directions {
		count := 1
		for i := 1; i < 4; i++ {
			c, r := col+d[0]*i, row+d[1]*i
			if c < 0 || c >= g.Columns || r < 0 || r >= g.Rows || g.Board[c][r] != player {
				break
			}
			count++
		}
		for i := 1; i < 4; i++ {
			c, r := col-d[0]*i, row-d[1]*i
			if c < 0 || c >= g.Columns || r < 0 || r >= g.Rows || g.Board[c][r] != player {
				break
			}
			count++
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) checkDraw() bool {
	for col := 0; col < g.Columns; col++ {
		for row := 0; row < g.Rows; row++ {
			if g.Board[col][row] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) scanForWinner() int {
	for col := 0; col < g.Columns; col++ {
		for row := 0; row < g.Rows; row++ {
			if g.Board[col][row] == 0 {
				continue
			}
			if g.checkWin(col, row) {
				return g.Board[col][row]
			}
		}
	}
	return 0
}
