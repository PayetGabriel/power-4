package game // package "game" contenant la logique du Puissance 4

import "encoding/json" // pour encoder/décoder JSON

// Game représente l'état du jeu Puissance 4
type Game struct {
	Board     [][]string `json:"board"`     // Plateau : matrice de strings ("empty", "red", "yellow")
	Turn      string     `json:"turn"`      // Joueur actuel : "red" ou "yellow"
	Winner    string     `json:"winner"`    // Gagnant : "", "red", "yellow" ou "draw"
	GameOver  bool       `json:"gameOver"`  // Indique si le jeu est terminé
	Mode      string     `json:"mode"`      // Mode de difficulté : easy, medium, normal, hard
	WinLength int        `json:"winLength"` // Nombre de jetons alignés nécessaire pour gagner
}

// MoveRequest représente une demande de coup
type MoveRequest struct {
	Column int `json:"column"` // colonne choisie par le joueur
}

// MoveResponse représente la réponse après un coup
type MoveResponse struct {
	Success bool   `json:"success"` // true si le coup a été joué
	Game    *Game  `json:"game"`    // état du jeu après le coup
	Message string `json:"message"` // message à afficher au joueur
}

// NewGame crée une nouvelle partie selon le mode
func NewGame(mode string) *Game {
	var rows, cols, winLength int // variables pour la configuration du plateau

	switch mode { // paramètres selon la difficulté
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

	board := make([][]string, rows) // création du plateau vide
	for i := 0; i < rows; i++ {
		row := make([]string, cols) // création d'une ligne vide
		for j := 0; j < cols; j++ {
			row[j] = "empty" // toutes les cases initialement vides
		}
		board[i] = row
	}

	return &Game{ // retourne la nouvelle partie
		Board:     board,
		Turn:      "red", // rouge commence
		Winner:    "",    // pas de gagnant
		GameOver:  false, // partie non terminée
		Mode:      mode,  // mode choisi
		WinLength: winLength,
	}
}

// MakeMove effectue un coup dans la colonne spécifiée
func (g *Game) MakeMove(col int) bool {
	if g.GameOver { // si le jeu est déjà terminé
		return false
	}

	if col < 0 || col >= len(g.Board[0]) { // colonne invalide
		return false
	}

	for row := len(g.Board) - 1; row >= 0; row-- { // recherche la première case vide de bas en haut
		if g.Board[row][col] == "empty" { // si case vide
			g.Board[row][col] = g.Turn // place le jeton du joueur actuel

			if g.checkWin() { // vérifie si ce coup fait gagner
				g.Winner = g.Turn
				g.GameOver = true
			} else if g.isBoardFull() { // sinon, vérifie si plateau plein (match nul)
				g.Winner = "draw"
				g.GameOver = true
			} else {
				g.switchPlayer() // change le joueur si partie continue
			}
			return true // coup joué
		}
	}

	return false // colonne pleine
}

// switchPlayer change le joueur actuel
func (g *Game) switchPlayer() {
	if g.Turn == "red" {
		g.Turn = "yellow" // passe de rouge à jaune
	} else {
		g.Turn = "red" // passe de jaune à rouge
	}
}

// checkWin vérifie s'il y a un gagnant
func (g *Game) checkWin() bool {
	for row := 0; row < len(g.Board); row++ { // parcours chaque ligne
		for col := 0; col < len(g.Board[0]); col++ { // parcours chaque colonne
			if g.Board[row][col] != "empty" { // ignore cases vides
				if g.checkDirection(row, col, 0, 1) || // horizontal
					g.checkDirection(row, col, 1, 0) || // vertical
					g.checkDirection(row, col, 1, 1) || // diagonale bas-droite
					g.checkDirection(row, col, 1, -1) { // diagonale bas-gauche
					return true
				}
			}
		}
	}
	return false // pas de victoire
}

// checkDirection vérifie s'il y a WinLength jetons alignés dans une direction donnée
func (g *Game) checkDirection(startRow, startCol, deltaRow, deltaCol int) bool {
	player := g.Board[startRow][startCol] // récupère le joueur à l'origine
	if player == "empty" {
		return false
	}

	count := 1 // compteur de jetons alignés

	row, col := startRow+deltaRow, startCol+deltaCol
	for count < g.WinLength && row >= 0 && row < len(g.Board) && col >= 0 && col < len(g.Board[0]) && g.Board[row][col] == player {
		count++         // incrémente le compteur si même joueur
		row += deltaRow // avance dans la direction
		col += deltaCol
	}

	row, col = startRow-deltaRow, startCol-deltaCol
	for count < g.WinLength && row >= 0 && row < len(g.Board) && col >= 0 && col < len(g.Board[0]) && g.Board[row][col] == player {
		count++ // compte jetons dans la direction opposée
		row -= deltaRow
		col -= deltaCol
	}

	return count >= g.WinLength // true si assez de jetons alignés
}

// isBoardFull vérifie si le plateau est plein (match nul)
func (g *Game) isBoardFull() bool {
	for col := 0; col < len(g.Board[0]); col++ { // parcours première ligne
		if g.Board[0][col] == "empty" { // si une case vide, pas full
			return false
		}
	}
	return true // toutes les cases remplies
}

// ToJSON convertit le jeu en JSON
func (g *Game) ToJSON() ([]byte, error) {
	return json.Marshal(g) // encode le struct Game en JSON
}

// Reset remet le jeu à zéro en conservant le mode
func (g *Game) Reset() {
	*g = *NewGame(g.Mode) // recrée une nouvelle partie avec le même mode
}
