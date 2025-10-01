package game

import "encoding/json"

// Game représente l'état du jeu Puissance 4
type Game struct {
	Board    [][]string `json:"board"`    // Plateau de jeu : "empty", "red", "yellow"
	Turn     string     `json:"turn"`     // Tour actuel : "red" ou "yellow"
	Winner   string     `json:"winner"`   // Gagnant : "", "red", "yellow", ou "draw"
	GameOver bool       `json:"gameOver"` // Indique si le jeu est terminé
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

// NewGame crée une nouvelle partie
func NewGame() *Game {
	board := make([][]string, ROWS)
	for i := 0; i < ROWS; i++ {
		row := make([]string, COLS)
		for j := 0; j < COLS; j++ {
			row[j] = "empty"
		}
		board[i] = row
	}
	return &Game{
		Board:    board,
		Turn:     "red",
		Winner:   "",
		GameOver: false,
	}
}

// Dimensions du plateau Puissance 4 classique
const ROWS = 6
const COLS = 7

// MakeMove effectue un coup dans la colonne spécifiée
func (g *Game) MakeMove(col int) bool {
	// Vérifier si le jeu est terminé
	if g.GameOver {
		return false
	}

	// Vérifier si la colonne est valide
	if col < 0 || col >= COLS {
		return false
	}

	// Trouver la première case vide dans la colonne (en partant du bas)
	for row := ROWS - 1; row >= 0; row-- {
		if g.Board[row][col] == "empty" {
			// Placer le jeton
			g.Board[row][col] = g.Turn

			// Vérifier la victoire
			if win, tokens := g.checkWin(); win {
				g.Winner = g.Turn
				g.GameOver = true
				// Enregistrer les jetons gagnants pour le frontend
				_ = tokens // Remplacez ceci par une logique pour transmettre les tokens au frontend
			} else if g.isBoardFull() {
				g.Winner = "draw"
				g.GameOver = true
			} else {
				// Changer de joueur
				g.switchPlayer()
			}
			return true
		}
	}

	// Colonne pleine
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

// checkWin vérifie s'il y a un gagnant et retourne les coordonnées des jetons alignés
func (g *Game) checkWin() (bool, [4][2]int) {
	var winningTokens [4][2]int

	// Vérifier toutes les positions possibles
	for row := 0; row < ROWS; row++ {
		for col := 0; col < COLS; col++ {
			if g.Board[row][col] != "empty" {
				// Vérifier les 4 directions : horizontal, vertical, diagonal /, diagonal \
				if g.checkDirectionWithTokens(row, col, 0, 1, &winningTokens) || // horizontal
					g.checkDirectionWithTokens(row, col, 1, 0, &winningTokens) || // vertical
					g.checkDirectionWithTokens(row, col, 1, 1, &winningTokens) || // diagonal \
					g.checkDirectionWithTokens(row, col, 1, -1, &winningTokens) { // diagonal /
					return true, winningTokens
				}
			}
		}
	}
	return false, winningTokens
}

// checkDirectionWithTokens vérifie s'il y a 4 jetons alignés dans une direction donnée et enregistre les coordonnées
func (g *Game) checkDirectionWithTokens(startRow, startCol, deltaRow, deltaCol int, tokens *[4][2]int) bool {
	player := g.Board[startRow][startCol]
	if player == "empty" {
		return false
	}

	count := 1 // Compter le jeton de départ
	(*tokens)[0] = [2]int{startRow, startCol}

	// Vérifier dans une direction
	row, col := startRow+deltaRow, startCol+deltaCol
	for count < 4 && row >= 0 && row < ROWS && col >= 0 && col < COLS && g.Board[row][col] == player {
		(*tokens)[count] = [2]int{row, col}
		count++
		row += deltaRow
		col += deltaCol
	}

	// Vérifier dans la direction opposée
	row, col = startRow-deltaRow, startCol-deltaCol
	for count < 4 && row >= 0 && row < ROWS && col >= 0 && col < COLS && g.Board[row][col] == player {
		(*tokens)[count] = [2]int{row, col}
		count++
		row -= deltaRow
		col -= deltaCol
	}

	return count >= 4
}

// isBoardFull vérifie si le plateau est plein (match nul)
func (g *Game) isBoardFull() bool {
	for col := 0; col < COLS; col++ {
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

// Reset remet le jeu à zéro
func (g *Game) Reset() {
	*g = *NewGame()
}
