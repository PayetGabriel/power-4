package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"power-4/src/game"
)

var tmpl = template.Must(template.ParseFiles("src/templates/index.html"))

var g = game.NewGame()

func handler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serveur démarré sur http://localhost:%d\n", port)

	// route pour le HTML
	http.HandleFunc("/", handler)

	// route pour les fichiers statiques (CSS et JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

	log.Fatal(http.Serve(listener, nil))

	
}
