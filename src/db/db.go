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

	// ⚠️ Le nom de base "power-4 db" contient un espace,
	// on doit l’entourer avec des backticks (``) dans le DSN,
	// OU mieux : la renommer en "power_4_db" dans phpMyAdmin.
	// En attendant, on l’écrit ainsi :
	dsn := "root:@tcp(127.0.0.1:3306)/power_4_db"

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Erreur de connexion à MySQL :", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Impossible de joindre la base :", err)
	}

	fmt.Println("✅ Connecté à MySQL (power_4_db)")
}

// UserExists retourne true si un utilisateur avec ce nom existe déjà
func UserExists(username string) (bool, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
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
	stmt := "INSERT INTO user (username, password) VALUES (?, ?)"
	res, err := DB.Exec(stmt, username, hashedPassword)
	if err != nil {
		log.Printf("CreateUser SQL error: %v", err)
		return err
	}

	// Log insert result pour débogage
	if id, err2 := res.LastInsertId(); err2 == nil {
		if rows, err3 := res.RowsAffected(); err3 == nil {
			log.Printf("CreateUser: LastInsertId=%d RowsAffected=%d", id, rows)
		} else {
			log.Printf("CreateUser: LastInsertId=%d RowsAffected(err)=%v", id, err3)
		}
	} else {
		log.Printf("CreateUser: LastInsertId(err)=%v", err2)
	}
	return nil
}
