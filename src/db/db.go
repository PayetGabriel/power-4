package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initialise la connexion MySQL
func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/power4" // adapte si ton mot de passe MySQL n’est pas vide
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à MySQL :", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Impossible de joindre la base :", err)
	}

	fmt.Println("✅ Connecté à MySQL")
}

// UserExists retourne true si un utilisateur avec ce nom existe déjà
func UserExists(username string) (bool, error) {
	var id int
	err := DB.QueryRow("SELECT Id FROM utilisateur WHERE Nom = ?", username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateUser insère un nouvel utilisateur (nom + mot de passe hashé)
func CreateUser(username, hashedPassword string) error {
	stmt := "INSERT INTO utilisateur (Nom, Mdp) VALUES (?, ?)"
	_, err := DB.Exec(stmt, username, hashedPassword)
	return err
}
