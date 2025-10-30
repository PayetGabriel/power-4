package main // déclaration du package principal

import (
	"encoding/json" // pour encoder/décoder JSON
	"fmt"           // pour le formatage de chaînes et print
	"html/template" // pour générer les pages HTML dynamiquement
	"log"           // pour le logging
	"net"           // pour écouter sur un port TCP
	"net/http"      // pour créer un serveur HTTP

	"power-4/src/game" // import du module de logique du jeu
)

var g *game.Game // variable globale pour stocker la partie en cours

// handler pour la page principale du jeu
func gameHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode") // récupère le paramètre "mode" depuis l'URL
	if mode == "" {                   // si aucun mode n'est fourni
		mode = "normal" // mode par défaut
	}

	g = game.NewGame(mode) // initialisation d'une nouvelle partie selon la difficulté

	tmpl := template.Must(template.ParseFiles("templates/index.html")) // parse le fichier HTML, panique si erreur
	err := tmpl.Execute(w, nil)                                        // exécute le template et écrit la réponse HTTP (ici, pas de données dynamiques)
	if err != nil {                                                    // si erreur lors de l'exécution
		http.Error(w, err.Error(), http.StatusInternalServerError) // renvoie une erreur 500
	}
}

// handler pour la page du menu
func menuHandler(w http.ResponseWriter, r *http.Request) {
	tmplMenu := template.Must(template.ParseFiles("templates/indexMenu.html")) // parse le template menu
	err := tmplMenu.Execute(w, nil)                                            // exécute et envoie la page
	if err != nil {                                                            // si erreur
		http.Error(w, err.Error(), http.StatusInternalServerError) // erreur 500
	}
}

// handler pour choisir la difficulté
func difficultyHandler(w http.ResponseWriter, r *http.Request) {
	tmplDiff := template.Must(template.ParseFiles("templates/difficulty.html")) // parse le template difficulté
	err := tmplDiff.Execute(w, nil)                                             // exécute et envoie la page
	if err != nil {                                                             // si erreur
		http.Error(w, err.Error(), http.StatusInternalServerError) // erreur 500
	}
}

// API pour récupérer l'état du jeu
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // indique que la réponse sera du JSON
	json.NewEncoder(w).Encode(g)                       // encode la partie actuelle en JSON et l'envoie
}

// API pour jouer un coup
func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // on n'accepte que les requêtes POST
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed) // sinon erreur 405
		return
	}

	var moveReq game.MoveRequest                    // structure pour recevoir le coup du joueur
	err := json.NewDecoder(r.Body).Decode(&moveReq) // décode le JSON reçu dans moveReq
	if err != nil {                                 // si JSON invalide
		http.Error(w, "JSON invalide", http.StatusBadRequest) // erreur 400
		return
	}

	success := g.MakeMove(moveReq.Column) // joue le coup dans la colonne demandée
	response := game.MoveResponse{        // prépare la réponse à renvoyer
		Success: success, // succès du coup
		Game:    g,       // état du jeu mis à jour
	}

	// gestion des messages selon succès et fin de partie
	if !success { // si coup invalide
		if g.GameOver { // si le jeu est terminé
			response.Message = "Le jeu est terminé" // message fin de partie
		} else {
			response.Message = "Coup invalide - colonne pleine ou incorrecte" // message coup impossible
		}
	} else if g.GameOver { // si coup valide mais fin de partie
		if g.Winner == "draw" { // match nul
			response.Message = "Match nul !"
		} else { // un joueur a gagné
			response.Message = fmt.Sprintf("Joueur %s gagne !", g.Winner) // message gagnant
		}
	}

	w.Header().Set("Content-Type", "application/json") // réponse en JSON
	json.NewEncoder(w).Encode(response)                // encode et envoie
}

// API pour réinitialiser le jeu
func resetGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // seulement POST
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	mode := r.URL.Query().Get("mode") // récupère le mode depuis l'URL
	if mode == "" {                   // mode par défaut si vide
		mode = "normal"
	}

	g = game.NewGame(mode) // recrée une nouvelle partie

	w.Header().Set("Content-Type", "application/json") // réponse JSON
	json.NewEncoder(w).Encode(g)                       // encode et envoie la nouvelle partie
}

// page de résultat
func resultHandler(w http.ResponseWriter, r *http.Request) {
	switch g.Winner { // selon le gagnant
	case "red": // si rouge
		t, _ := template.ParseFiles("templates/win.html")  // parse template victoire
		t.Execute(w, map[string]string{"Winner": "Rouge"}) // exécute avec variable Winner
	case "yellow": // si jaune
		t, _ := template.ParseFiles("templates/win.html")
		t.Execute(w, map[string]string{"Winner": "Jaune"})
	case "draw": // match nul
		t, _ := template.ParseFiles("templates/draw.html")
		t.Execute(w, nil)
	default: // si pas de gagnant encore
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirection vers menu principal
	}
}

// replay : reset + redirection en conservant le mode (depuis query param)
func replayHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")           // récupère mode depuis query param
	if mode == "" && g != nil && g.Mode != "" { // si vide mais game existe
		mode = g.Mode // conserve le mode actuel
	}
	if mode == "" { // sinon mode par défaut
		mode = "normal"
	}

	// recréer clairement la partie dans le mode demandé
	g = game.NewGame(mode)

	// rediriger vers /game?mode=... pour que le frontend récupère le bon plateau
	http.Redirect(w, r, "/game?mode="+mode, http.StatusSeeOther) // w = réponse HTTP, r = requête reçue
}

func main() {
	listener, err := net.Listen("tcp", ":0") // écoute sur un port TCP libre
	if err != nil {                          // si erreur
		panic(err) // stoppe le programme
	}

	port := listener.Addr().(*net.TCPAddr).Port                   // récupère le port réel choisi
	fmt.Printf("Serveur démarré sur http://localhost:%d\n", port) // affiche le lien d'accès

	// routes principales
	http.HandleFunc("/", menuHandler)                 // route menu
	http.HandleFunc("/difficulty", difficultyHandler) // route difficulté
	http.HandleFunc("/game", gameHandler)             // route jeu principal

	// routes API
	http.HandleFunc("/api/game-state", gameStateHandler) // récupérer état du jeu
	http.HandleFunc("/api/make-move", makeMoveHandler)   // jouer un coup
	http.HandleFunc("/api/reset-game", resetGameHandler) // réinitialiser la partie

	http.HandleFunc("/result", resultHandler) // page résultat
	http.HandleFunc("/replay", replayHandler) // rejouer partie

	// servir fichiers statiques (CSS, JS, images) depuis dossier /static/
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.Serve(listener, nil)) // démarre le serveur et log fatal si erreur
}
