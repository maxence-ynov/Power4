package main

type Game struct {
	Board   [][]int // 0 = empty, 1 = player1, 2 = player2
	Rows    int
	Columns int
	Turn    int // 1 or 2
	Winner  int // 0 = none, 1 or 2
	Draw    bool
}

func NewGame() *Game {
	rows, cols := 6, 7
	board := make([][]int, cols)
	for i := range board {
		board[i] = make([]int, rows)
	}
	return &Game{
		Board:   board,
		Rows:    rows,
		Columns: cols,
		Turn:    1,
		Winner:  0,
		Draw:    false,
	}
}

// Play places a token in the chosen column
func (g *Game) Play(col int) {
	if col < 0 || col >= g.Columns || g.Winner != 0 || g.Draw {
		return
	}
	// Find the lowest empty cell in the column
	for row := g.Rows - 1; row >= 0; row-- {
		if g.Board[col][row] == 0 {
			g.Board[col][row] = g.Turn
			if g.checkWin(col, row) {
				g.Winner = g.Turn
			} else if g.checkDraw() {
				g.Draw = true
			} else {
				g.Turn = 3 - g.Turn // Switch player
			}
			break
		}
	}
}

// Check for win condition
func (g *Game) checkWin(col, row int) bool {
	player := g.Board[col][row]
	directions := [][2]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
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

// Check for draw
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
