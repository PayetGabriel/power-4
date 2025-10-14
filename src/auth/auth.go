package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // haché
}

// On garde les utilisateurs dans un fichier JSON (ex : src/auth/users.json)
const userFile = "src/auth/users.json"

func hashPassword(pw string) string {
	hash := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(hash[:])
}

// Crée un nouveau compte
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "src/templates/signup.html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	users := loadUsers()
	for _, u := range users {
		if u.Username == username {
			http.Error(w, "Nom déjà utilisé", http.StatusConflict)
			return
		}
	}

	users = append(users, User{Username: username, Password: hashPassword(password)})
	saveUsers(users)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Connexion utilisateur
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "src/templates/login.html")
		return
	}

	username := r.FormValue("username")
	password := hashPassword(r.FormValue("password"))

	users := loadUsers()
	for _, u := range users {
		if u.Username == username && u.Password == password {
			http.SetCookie(w, &http.Cookie{Name: "user", Value: username, Path: "/"})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
}

// Utilitaire pour lire/écrire les utilisateurs
func loadUsers() []User {
	var users []User
	data, err := os.ReadFile(userFile)
	if err == nil {
		json.Unmarshal(data, &users)
	}
	return users
}

func saveUsers(users []User) {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(userFile, data, 0644)
}

// Récupérer l’utilisateur connecté
func GetLoggedUser(r *http.Request) (string, error) {
	c, err := r.Cookie("user")
	if err != nil {
		return "", errors.New("non connecté")
	}
	return c.Value, nil
}
