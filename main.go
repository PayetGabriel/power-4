package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	// On demande au système d'attribuer un port libre (":0")
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	// On récupère le port attribué
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serveur démarré sur le port %d...\n", port)

	// On lance le serveur
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur mon serveur Go !")
	})

	http.Serve(listener, nil)
}
