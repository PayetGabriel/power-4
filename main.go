package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"power-4/src/game" // adapte selon ton go.mod (par ex. "power4/src/game")
)

var tmpl = template.Must(template.ParseFiles("src/templates/index.html"))

// On instancie une partie de Puissance 4
var g = game.NewGame()

func handler(w http.ResponseWriter, r *http.Request) {
	// On passe le plateau au template pour affichage
	err := tmpl.Execute(w, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// On demande un port libre au système
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	// On récupère le port attribué
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serveur démarré sur http://localhost:%d\n", port)

	// On configure la route principale
	http.HandleFunc("/", handler)

	// On lance le serveur
	log.Fatal(http.Serve(listener, nil))
}
