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

var tmpl = template.Must(template.ParseFiles("src/templates/index.html"))

var g = game.NewGame()

// handler pour la page principale
func handler(w http.ResponseWriter, r *http.Request) {
	// afficher la page principale (index.html) — le JS fera les requêtes /api/game-state
	err := tmpl.Execute(w, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// gameStateHandler retourne l'état actuel du jeu en JSON
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// makeMoveHandler traite un coup de joueur
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

// resetGameHandler remet le jeu à zéro
func resetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	g.Reset()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// resultHandler affiche la page de résultat (victoire / match nul)
func resultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner {
	case "red":
		t, err := template.ParseFiles("src/templates/win.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// on fournit le nom lisible du gagnant
		data := map[string]string{"Winner": "Rouge"}
		t.Execute(w, data)
	case "yellow":
		t, err := template.ParseFiles("src/templates/win.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[string]string{"Winner": "Jaune"}
		t.Execute(w, data)
	case "draw":
		t, err := template.ParseFiles("src/templates/draw.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	default:
		// pas de résultat -> redirige vers le jeu
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// replayHandler réinitialise le jeu puis renvoie à la page d'accueil
func replayHandler(w http.ResponseWriter, r *http.Request) {
	g.Reset()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serveur démarré sur http://localhost:%d\n", port)

	// Route pour la page principale
	http.HandleFunc("/", handler)

	// Routes API
	http.HandleFunc("/api/game-state", gameStateHandler)
	http.HandleFunc("/api/make-move", makeMoveHandler)
	http.HandleFunc("/api/reset-game", resetGameHandler)

	// Route résultat
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/replay", replayHandler)

	// Route pour les fichiers statiques (CSS et JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

	log.Fatal(http.Serve(listener, nil))
}
