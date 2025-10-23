package main

type Game struct {
	Board           [][]int
	Rows            int
	Columns         int
	Turn            int
	Winner          int
	Draw            bool
	GravityUp       bool
	TurnCount       int
	LastCol         int
	LastRow         int
	LastGravityUsed bool
	Player1Name     string
	Player2Name     string
	Level           string
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
		Board:           board,
		Rows:            rows,
		Columns:         cols,
		Turn:            1,
		Winner:          0,
		Draw:            false,
		TurnCount:       0,
		LastCol:         -1,
		LastRow:         -1,
		LastGravityUsed: false,
		Player1Name:     player1Name,
		Player2Name:     player2Name,
		Level:           level,
	}
}

func (g *Game) Play(col int) {
	if col < 0 || col >= g.Columns || g.Winner != 0 || g.Draw {
		return
	}

	// Flip gravity ahead of the move when the upcoming turn reaches the threshold so the drop animation matches the new direction.
	moveGravity := g.GravityUp
	nextTurn := g.TurnCount + 1
	shouldToggle := nextTurn%5 == 0
	prevGravity := g.GravityUp
	if shouldToggle {
		g.GravityUp = !g.GravityUp
	}

	var row int
	row = -1
	if g.GravityUp {
		for r := 0; r < g.Rows; r++ {
			if g.Board[col][r] == 0 {
				g.Board[col][r] = g.Turn
				row = r
				break
			}
		}
	} else {
		for r := g.Rows - 1; r >= 0; r-- {
			if g.Board[col][r] == 0 {
				g.Board[col][r] = g.Turn
				row = r
				break
			}
		}
	}

	if row == -1 {
		if shouldToggle {
			g.GravityUp = prevGravity
		}
		return
	}

	g.Board[col][row] = g.Turn

	// Apply gravity once so every piece settles in the active direction and the frontend can mirror it.
	row = g.applyGravity(col, row, g.GravityUp)
	g.LastCol = col
	g.LastRow = row
	g.LastGravityUsed = moveGravity
	g.TurnCount++

	if g.checkWin(col, row) {
		g.Winner = g.Turn
		return
	}

	if g.checkDraw() {
		g.Draw = true
		return
	}

	g.Turn = 3 - g.Turn
}

func (g *Game) applyGravity(lastCol, lastRow int, gravityUp bool) int {
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

		count := len(stack)
		for i, entry := range stack {
			var newRow int
			if gravityUp {
				newRow = i
			} else {
				// Preserve bottom-up order when packing tokens under normal gravity.
				newRow = g.Rows - count + i
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
