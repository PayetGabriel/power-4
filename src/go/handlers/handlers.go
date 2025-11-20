package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"power-4/src/go/db"
	"power-4/src/go/game"

	"golang.org/x/crypto/bcrypt"
)

var (
	tmplGame       = template.Must(template.ParseFiles("templates/html/index.html"))
	tmplMenu       = template.Must(template.ParseFiles("templates/html/indexMenu.html"))
	tmplDifficulty = template.Must(template.ParseFiles("templates/html/difficulty.html"))
	g              = game.NewGame("normal")
)

// --- PAGE D'ACCUEIL (MENU) ---
func MenuHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmplMenu.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- PAGE DE CHOIX DE DIFFICULTÉ ---
func DifficultyHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmplDifficulty.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- PAGE DE JEU ---
func GameHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode != "" {
		log.Printf("🎮 Nouvelle partie en mode: %s", mode)
		g = game.NewGame(mode)
	}
	if err := tmplGame.Execute(w, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- INSCRIPTION ---
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.Must(template.ParseFiles("templates/html/signup.html"))
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
		http.Redirect(w, r, "/difficulty", http.StatusSeeOther)
		return
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// --- CONNEXION UTILISATEUR ---
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.Must(template.ParseFiles("templates/html/login.html"))
		t.Execute(w, nil)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
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
		http.Redirect(w, r, "/difficulty", http.StatusSeeOther)
		return
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// --- ÉTAT DU JEU ---
func GameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// --- TRAITEMENT D'UN COUP ---
func MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
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
	response := game.MoveResponse{Success: success, Game: g}
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
func ResetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	g.Reset()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// --- PAGE DE RÉSULTAT ---
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner {
	case "red", "yellow":
		t, err := template.ParseFiles("templates/html/win.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		color := map[string]string{"red": "Rouge", "yellow": "Jaune"}[g.Winner]
		t.Execute(w, map[string]string{"Winner": color})
	case "draw":
		t, err := template.ParseFiles("templates/html/draw.html")
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
func ReplayHandler(w http.ResponseWriter, r *http.Request) {
	g.Reset()
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

type MimeFileServer struct{ Handler http.Handler }

func (m *MimeFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	m.Handler.ServeHTTP(w, r)
}
