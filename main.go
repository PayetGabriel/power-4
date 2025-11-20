package main

import (
	"power-4/src/go/db"
	"power-4/src/go/server"
)

// Entrypoint
func main() {
	db.InitDB()
	defer db.DB.Close()
	server.StartServer()
}
