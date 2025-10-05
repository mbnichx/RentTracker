package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("super_secret_key") //  Use env variable in production

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// generateJWT returns a signed JWT token
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// POST /register
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Print(err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		log.Print(creds)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		creds.UserPassword = string(hashedPassword)

		err = CreateUser(db, &creds)
		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		token, _ := generateJWT(int(creds.UserID))

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// POST /login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var userId int
		var email int
		var hashedPassword string

		err := db.QueryRow("SELECT userId, email, password FROM users WHERE email = ?", creds.Email).Scan(&userId, &email, &hashedPassword)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := generateJWT(userId)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
