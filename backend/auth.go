/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file provides authentication logic for login and JWT token generation.
// Includes the Credentials struct, generateJWT helper, and LoginHandler for
// validating user credentials and issuing tokens. Uses bcrypt for password
// verification and JWT for session management.

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
// generateJWT creates and signs a JWT token for the given user ID.
// Returns the token string and error if signing fails.
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// POST /login
// LoginHandler returns an HTTP handler for user login.
// Accepts a JSON body with credentials, verifies password, and responds with a JWT token.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		log.Print("Attempting to log in " + creds.Email + " ...")

		var userId int
		var hashedPassword string

		err := db.QueryRow("SELECT userId, userPasswordHash FROM users WHERE userEmail = ?", creds.Email).Scan(&userId, &hashedPassword)
		if err != nil {
			log.Print(err)
			http.Error(w, "User does not exist", http.StatusUnauthorized)
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
