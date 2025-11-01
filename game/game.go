package game

import "encoding/json"

// Game représente l'état du jeu Puissance 4
type Game struct {
	Board     [][]string `json:"board"`
	Turn      string     `json:"turn"`
	Winner    string     `json:"winner"`
	GameOver  bool       `json:"gameOver"`
	Mode      string     `json:"mode"`
	WinLength int        `json:"winLength"`
}

// MoveRequest représente une demande de coup
type MoveRequest struct {
	Column int `json:"column"`
}

// MoveResponse représente la réponse après un coup
type MoveResponse struct {
	Success bool   `json:"success"`
	Game    *Game  `json:"game"`
	Message string `json:"message"`
}

// NewGame crée une nouvelle partie selon le mode
func NewGame(mode string) *Game {
	var rows, cols, winLength int

	switch mode {
	case "easy":
		rows, cols, winLength = 6, 7, 3
	case "medium":
		rows, cols, winLength = 6, 9, 5
	case "normal":
		rows, cols, winLength = 6, 7, 4
	case "hard":
		rows, cols, winLength = 7, 8, 7
	default:
		rows, cols, winLength = 6, 7, 4
	}

	board := make([][]string, rows)
	for i := 0; i < rows; i++ {
		row := make([]string, cols)
		for j := 0; j < cols; j++ {
			row[j] = "empty"
		}
		board[i] = row
	}

	return &Game{
		Board:     board,
		Turn:      "red",
		Winner:    "",
		GameOver:  false,
		Mode:      mode,
		WinLength: winLength,
	}
}

// MakeMove effectue un coup dans la colonne spécifiée
func (g *Game) MakeMove(col int) bool {
	if g.GameOver {
		return false
	}

	if col < 0 || col >= len(g.Board[0]) {
		return false
	}

	for row := len(g.Board) - 1; row >= 0; row-- {
		if g.Board[row][col] == "empty" {
			g.Board[row][col] = g.Turn

			if g.checkWin() {
				g.Winner = g.Turn
				g.GameOver = true
			} else if g.isBoardFull() {
				g.Winner = "draw"
				g.GameOver = true
			} else {
				g.switchPlayer()
			}
			return true
		}
	}

	return false
}

// switchPlayer change le joueur actuel
func (g *Game) switchPlayer() {
	if g.Turn == "red" {
		g.Turn = "yellow"
	} else {
		g.Turn = "red"
	}
}

// checkWin vérifie s'il y a un gagnant
func (g *Game) checkWin() bool {
	for row := 0; row < len(g.Board); row++ {
		for col := 0; col < len(g.Board[0]); col++ {
			if g.Board[row][col] != "empty" {
				if g.checkDirection(row, col, 0, 1) ||
					g.checkDirection(row, col, 1, 0) ||
					g.checkDirection(row, col, 1, 1) ||
					g.checkDirection(row, col, 1, -1) {
					return true
				}
			}
		}
	}
	return false
}

// checkDirection vérifie s'il y a N jetons alignés dans une direction donnée
func (g *Game) checkDirection(startRow, startCol, deltaRow, deltaCol int) bool {
	player := g.Board[startRow][startCol]
	if player == "empty" {
		return false
	}

	count := 1

	row, col := startRow+deltaRow, startCol+deltaCol
	for count < g.WinLength && row >= 0 && row < len(g.Board) && col >= 0 && col < len(g.Board[0]) && g.Board[row][col] == player {
		count++
		row += deltaRow
		col += deltaCol
	}

	row, col = startRow-deltaRow, startCol-deltaCol
	for count < g.WinLength && row >= 0 && row < len(g.Board) && col >= 0 && col < len(g.Board[0]) && g.Board[row][col] == player {
		count++
		row -= deltaRow
		col -= deltaCol
	}

	return count >= g.WinLength
}

// isBoardFull vérifie si le plateau est plein (match nul)
func (g *Game) isBoardFull() bool {
	for col := 0; col < len(g.Board[0]); col++ {
		if g.Board[0][col] == "empty" {
			return false
		}
	}
	return true
}

// ToJSON convertit le jeu en JSON
func (g *Game) ToJSON() ([]byte, error) {
	return json.Marshal(g)
}

// Reset remet le jeu à zéro en conservant le mode
func (g *Game) Reset() {
	*g = *NewGame(g.Mode)
}
