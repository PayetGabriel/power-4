package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"power-4/src/go/handlers"
)

// StartServer configure les routes et démarre le serveur HTTP.
func StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute : %v", err)
	}
	fmt.Printf("🚀 Serveur démarré sur http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	// Routes
	http.HandleFunc("/", handlers.MenuHandler)
	http.HandleFunc("/menu", handlers.MenuHandler)
	http.HandleFunc("/difficulty", handlers.DifficultyHandler)
	http.HandleFunc("/game", handlers.GameHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/api/game-state", handlers.GameStateHandler)
	http.HandleFunc("/api/make-move", handlers.MakeMoveHandler)
	http.HandleFunc("/api/reset-game", handlers.ResetGameHandler)
	http.HandleFunc("/result", handlers.ResultHandler)
	http.HandleFunc("/replay", handlers.ReplayHandler)

	// Fichiers statiques
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", &handlers.MimeFileServer{Handler: fs}))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("src/js"))))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() { serverErr <- http.Serve(listener, nil) }()

	select {
	case err := <-serverErr:
		log.Fatalf("Erreur serveur : %v", err)
	case <-stop:
		fmt.Println("\n🛑 Serveur interrompu proprement.")
	}
}
