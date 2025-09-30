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
