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
	"power-4/db"
	"power-4/game"
	"syscall"

	"golang.org/x/crypto/bcrypt"
)

var (
	tmplGame       = template.Must(template.ParseFiles("web/templates/index.html"))
	tmplMenu       = template.Must(template.ParseFiles("web/templates/indexMenu.html"))
	tmplDifficulty = template.Must(template.ParseFiles("web/templates/difficulty.html"))
	g              = game.NewGame("normal")
)

// --- PAGE D'ACCUEIL (MENU) ---
func menuHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplMenu.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- PAGE DE CHOIX DE DIFFICULTÉ ---
func difficultyHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplDifficulty.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- PAGE DE JEU ---
func gameHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le mode depuis l'URL (ex: /game?mode=easy)
	mode := r.URL.Query().Get("mode")

	// Si un mode est spécifié, crée une nouvelle partie avec ce mode
	if mode != "" {
		log.Printf("🎮 Nouvelle partie en mode: %s", mode)
		g = game.NewGame(mode)
	}

	err := tmplGame.Execute(w, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- INSCRIPTION ---
func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.Must(template.ParseFiles("web/templates/signup.html"))
		t.Execute(w, nil)
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Veuillez remplir tous les champs", http.StatusBadRequest)
			return
		}

		exists, err := db.UserExists(username)
		if err != nil {
			log.Printf("Erreur UserExists: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Nom d'utilisateur déjà pris", http.StatusConflict)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Erreur bcrypt: %v", err)
			http.Error(w, "Erreur lors du hashage", http.StatusInternalServerError)
			return
		}

		if err := db.CreateUser(username, string(hashed)); err != nil {
			log.Printf("Erreur CreateUser: %v", err)
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Utilisateur créé : %s", username)
		// Redirection vers la page de difficulté après inscription
		http.Redirect(w, r, "/difficulty", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// --- CONNEXION UTILISATEUR ---
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.Must(template.ParseFiles("web/templates/login.html"))
		t.Execute(w, nil)
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Vérifie si l'utilisateur existe
		exists, err := db.UserExists(username)
		if err != nil {
			log.Printf("Erreur UserExists lors du login: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		if !exists {
			log.Printf("❌ Tentative de connexion : utilisateur '%s' non trouvé", username)
			http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
			return
		}

		// Vérifie le mot de passe
		var hashedPassword string
		err = db.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			log.Printf("Erreur récupération mot de passe: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Printf("❌ Mot de passe incorrect pour '%s'", username)
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		log.Printf("✅ Connexion réussie : %s", username)
		// Connexion réussie → redirection vers choix de difficulté
		http.Redirect(w, r, "/difficulty", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// --- ÉTAT DU JEU ---
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// --- TRAITEMENT D'UN COUP ---
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

// --- RÉINITIALISATION DU JEU ---
func resetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	g.Reset()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// --- PAGE DE RÉSULTAT ---
func resultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner {
	case "red", "yellow":
		t, err := template.ParseFiles("web/templates/win.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		color := map[string]string{"red": "Rouge", "yellow": "Jaune"}[g.Winner]
		t.Execute(w, map[string]string{"Winner": color})
	case "draw":
		t, err := template.ParseFiles("web/templates/draw.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// --- REJOUER UNE PARTIE ---
func replayHandler(w http.ResponseWriter, r *http.Request) {
	g.Reset()
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// --- MAIN ---
func main() {
	// 🔗 Connexion à la base MySQL
	db.InitDB()
	defer db.DB.Close()

	// 🔥 Lancement du serveur
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute : %v", err)
	}
	fmt.Printf("🚀 Serveur démarré sur http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	// --- ROUTES ---
	http.HandleFunc("/", menuHandler)                 // Page d'accueil = menu
	http.HandleFunc("/menu", menuHandler)             // Menu principal
	http.HandleFunc("/difficulty", difficultyHandler) // ✅ Choix de difficulté
	http.HandleFunc("/game", gameHandler)             // Grille de jeu (+ gestion du mode)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/api/game-state", gameStateHandler)
	http.HandleFunc("/api/make-move", makeMoveHandler)
	http.HandleFunc("/api/reset-game", resetGameHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/replay", replayHandler)

	// --- FICHIERS STATIQUES ---
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// --- GESTION DE L'ARRÊT PROPRE ---
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
