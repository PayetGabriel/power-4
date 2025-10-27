package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"power-4/src/game"
)

var g *game.Game

// handler pour la page principale du jeu
func gameHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "normal"
	}

	g = game.NewGame(mode) // initialisation selon difficulté

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil) // <--- PAS g
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handler pour la page du menu
func menuHandler(w http.ResponseWriter, r *http.Request) {
	tmplMenu := template.Must(template.ParseFiles("templates/indexMenu.html"))
	err := tmplMenu.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handler pour choisir la difficulté
func difficultyHandler(w http.ResponseWriter, r *http.Request) {
	tmplDiff := template.Must(template.ParseFiles("templates/difficulty.html"))
	err := tmplDiff.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// API pour récupérer l'état du jeu
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// API pour jouer un coup
func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var moveReq game.MoveRequest
	err := json.NewDecoder(r.Body).Decode(&moveReq)
	if err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	success := g.MakeMove(moveReq.Column)
	response := game.MoveResponse{
		Success: success,
		Game:    g,
	}

	if !success {
		if g.GameOver {
			response.Message = "Le jeu est terminé"
		} else {
			response.Message = "Coup invalide - colonne pleine ou incorrecte"
		}
	} else if g.GameOver {
		if g.Winner == "draw" {
			response.Message = "Match nul !"
		} else {
			response.Message = fmt.Sprintf("Joueur %s gagne !", g.Winner)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// API pour réinitialiser le jeu
func resetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "normal"
	}

	g = game.NewGame(mode)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// page de résultat
func resultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner {
	case "red":
		t, _ := template.ParseFiles("templates/win.html")
		t.Execute(w, map[string]string{"Winner": "Rouge"})
	case "yellow":
		t, _ := template.ParseFiles("templates/win.html")
		t.Execute(w, map[string]string{"Winner": "Jaune"})
	case "draw":
		t, _ := template.ParseFiles("templates/draw.html")
		t.Execute(w, nil)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// replay : reset + redirection
func replayHandler(w http.ResponseWriter, r *http.Request) {
	g.Reset()
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serveur démarré sur http://localhost:%d\n", port)

	// routes
	http.HandleFunc("/", menuHandler)
	http.HandleFunc("/difficulty", difficultyHandler)
	http.HandleFunc("/game", gameHandler)

	http.HandleFunc("/api/game-state", gameStateHandler)
	http.HandleFunc("/api/make-move", makeMoveHandler)
	http.HandleFunc("/api/reset-game", resetGameHandler)

	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/replay", replayHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.Serve(listener, nil))
}
