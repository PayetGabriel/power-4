package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"power-4/src/db"
	"power-4/src/game"
	"syscall"

	"golang.org/x/crypto/bcrypt"
)

var (
	tmpl = template.Must(template.ParseFiles("templates/index.html"))
	g    = game.NewGame()
)

// Page principale du jeu
func handler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Menu principal
func menuHandler(w http.ResponseWriter, r *http.Request) {
	tmplMenu := template.Must(template.ParseFiles("templates/indexMenu.html"))
	err := tmplMenu.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Affiche/traite le formulaire d'inscription
func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.Must(template.ParseFiles("templates/signup.html"))
		t.Execute(w, nil)
		return
	case http.MethodPost:
		log.Printf("register: request from %s method=%s", r.RemoteAddr, r.Method)
		// parse form
		if err := r.ParseForm(); err != nil {
			log.Println("register: parse form error:", err)
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Printf("register: form values username=%s (password hidden)", username)
		if username == "" || password == "" {
			http.Error(w, "Remplir tous les champs", http.StatusBadRequest)
			return
		}

		// check doublon
		exists, err := db.UserExists(username)
		if err != nil {
			log.Println("register: UserExists error:", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusConflict)
			return
		}

		// hash password
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("register: bcrypt error:", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		log.Println("register: creating user", username)
		if err := db.CreateUser(username, string(hashed)); err != nil {
			log.Println("register: CreateUser error:", err)
			http.Error(w, "Impossible de créer l'utilisateur", http.StatusInternalServerError)
			return
		}
		log.Println("register: user created", username)

		// Rediriger vers la page de login pour usage normal via navigateur
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}

// État du jeu en JSON
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// Traitement d’un coup
func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var moveReq game.MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
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

// Réinitialisation du jeu
func resetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	g.Reset()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// Affichage du résultat
func resultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner {
	case "red", "yellow":
		t, err := template.ParseFiles("templates/win.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		color := map[string]string{"red": "Rouge", "yellow": "Jaune"}[g.Winner]
		t.Execute(w, map[string]string{"Winner": color})
	case "draw":
		t, err := template.ParseFiles("templates/draw.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Rejouer une partie
func replayHandler(w http.ResponseWriter, r *http.Request) {
	g.Reset()
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func main() {
	db.InitDB()

	// Use a fixed port to avoid confusion with a random port (":0").
	// This makes it easier to access the server consistently during tests.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute : %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println("🎮Ctrl + Click sur le lien pour lancer")
	fmt.Printf("✅ Serveur démarré sur http://localhost:%d\n", port)
	fmt.Println("🛑Ctrl + C pour aretter le programme")

	// Routes
	http.HandleFunc("/", menuHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/game", handler)
	http.HandleFunc("/api/game-state", gameStateHandler)
	http.HandleFunc("/api/make-move", makeMoveHandler)
	http.HandleFunc("/api/reset-game", resetGameHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/replay", replayHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Gestion des signaux d'arrêt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- http.Serve(listener, nil)
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("Erreur serveur : %v", err)
	case <-stop:
		fmt.Println("\n🛑 Serveur interrompu proprement.")
	}
}
