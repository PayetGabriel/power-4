package game

type Game struct {
	Board [][]string // "empty", "p1", "p2"
	Turn  int        // 1 ou 2
}

func NewGame() *Game {
	board := make([][]string, 6)
	for i := 0; i < 6; i++ {
		row := make([]string, 7)
		for j := 0; j < 7; j++ {
			row[j] = "empty"
		}
		board[i] = row
	}
	return &Game{Board: board, Turn: 1}
}

// checkWin vérifie s’il y a un gagnant
func CheckWin() string {
	directions := [][]int{
		{1, 2, 3},                  // horizontal
		{COLS, 2 * COLS, 3 * COLS}, // vertical
		{COLS - 1, 2 * (COLS - 1), 3 * (COLS - 1)}, // diagonal /
		{COLS + 1, 2 * (COLS + 1), 3 * (COLS + 1)}, // diagonal \
	}

	for i := 0; i < len(cells); i++ {
		for _, deltas := range directions {
			sequence := KVoisins(i, deltas)

			// Vérifie "red"
			if len(sequence) == 3 && AllEqual(sequence, "red") && cells[i] == "red" {
				return "Red wins!"
			}

			// Vérifie "yellow"
			if len(sequence) == 3 && AllEqual(sequence, "yellow") && cells[i] == "yellow" {
				return "Yellow wins!"
			}
		}
	}
	return "" // Pas encore de gagnant
}

// allEqual vérifie si tous les éléments d’un tableau sont égaux à une valeur donnée
func AllEqual(arr []string, val string) bool {
	for _, v := range arr {
		if v != val {
			return false
		}
	}
	return true
}

// Dimensions du plateau (par ex. Puissance 4 classique : 6 lignes x 7 colonnes)
const ROWS = 6
const COLS = 7

// cells : plateau représenté comme un tableau 1D
// Chaque case contient "red", "yellow", ou "" si vide
var cells = make([]string, ROWS*COLS)

// KVoisins retourne les cases voisines en fonction des deltas
func KVoisins(index int, deltas []int) []string {
	voisins := []string{}
	for _, delta := range deltas {
		voisinIndex := index + delta
		if voisinIndex >= 0 && voisinIndex < len(cells) {
			voisins = append(voisins, cells[voisinIndex])
		}
	}
	return voisins
}

func (g *Game) DropToken(col int) (int, bool) {
	for row := ROWS - 1; row >= 0; row-- {
		if g.Board[row][col] == "empty" {
			if g.Turn == 1 {
				g.Board[row][col] = "p1"
				CheckWin()
			} else {
				g.Board[row][col] = "p2"
				CheckWin()
			}
			g.Turn = 3 - g.Turn // alterne entre 1 et 2
			return row, true
		}
	}
	return -1, false // colonne pleine
}
