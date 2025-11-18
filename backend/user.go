/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file implements user registration, retrieval, update, and deletion HTTP handlers
// and database helpers. Handlers: CreateUserHandler, GetUserHandler, GetCurrentUserHandler,
// UpdateUserHandler, DeleteUserHandler. DB helpers: CreateUser, GetAllUsers, GetUserByID,
// UpdateUser, DeleteUser. Passwords are hashed using bcrypt.

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// USERS
type User struct {
	UserID          int    `db:"userId" json:"userId"`
	UserFirstName   string `db:"userFirstName" json:"userFirstName"`
	UserLastName    string `db:"userLastName" json:"userLastName"`
	UserEmail       string `db:"userEmail" json:"userEmail"`
	UserPhoneNumber string `db:"userPhoneNumber" json:"userPhoneNumber"`
	UserPassword    string `db:"userPasswordHash" json:"userPassword"`
	UserRole        string `db:"userRole" json:"userRole"`
}

// == Handlers ========================================================================

// CreateUserHandler returns an HTTP handler for creating a new user.
// Accepts a JSON body with user details, hashes the password, and inserts the user into the database.
// Responds with the created user or error.
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			log.Print(err)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		u.UserPassword = string(hashedPassword)

		err = CreateUser(db, &u)
		if err != nil {
			log.Print(err)
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, u)
	}
}

// GetUserHandler returns an HTTP handler for retrieving user(s).
// If no userId is provided, returns all users. Otherwise, returns the user with the given ID.
func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		if idStr == "" {
			users, err := GetAllUsers(db)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, users)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid userId")
			return
		}
		u, err := GetUserByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondJSON(w, http.StatusOK, u)
	}
}

// GetCurrentUserHandler returns an HTTP handler for retrieving the current user.
// If no userId is provided, returns all users. Otherwise, returns the user with the given ID.
func GetCurrentUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		if idStr == "" {
			users, err := GetAllUsers(db)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, users)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid userId")
			return
		}
		u, err := GetUserByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondJSON(w, http.StatusOK, u)
	}
}

// UpdateUserHandler returns an HTTP handler for updating user details.
// Accepts a JSON body with updated user info and updates the user in the database.
func UpdateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if u.UserID == 0 {
			respondError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		if err := UpdateUser(db, &u); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DeleteUserHandler returns an HTTP handler for deleting a user by ID.
// Accepts a DELETE request and removes the user from the database.
func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid userId")
			return
		}
		if err := DeleteUser(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL helpers ========================================================================

// CreateUser inserts a new user into the database.
// Hashes the password and sets the UserID field on success.
// Returns an error if insertion fails.
func CreateUser(db *sql.DB, u *User) error {
	query := `INSERT INTO users (userFirstName, userLastName, userEmail, userPhoneNumber, userPasswordHash, userRole)
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, u.UserFirstName, u.UserLastName, u.UserEmail, u.UserPhoneNumber, u.UserPassword, u.UserRole)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.UserID = int(id)
	return nil
}

// UpdateUser updates an existing user's details in the database.
// Returns an error if the update fails.
func UpdateUser(db *sql.DB, u *User) error {
	_, err := db.Exec(`UPDATE users 
		SET userFirstName=?, userLastName=?, userEmail=?, userPhoneNumber=?, userRole=? 
		WHERE userId=?`,
		u.UserFirstName, u.UserLastName, u.UserEmail, u.UserPhoneNumber, u.UserRole, u.UserID)
	return err
}

// DeleteUser removes a user from the database by userId.
// Returns an error if deletion fails.
func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE userId=?`, id)
	return err
}

// READ
// GetAllUsers retrieves all users from the database.
// Returns a slice of User and an error if the query fails.
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT userId, userFirstName, userLastName, userEmail, userPhoneNumber, userRole FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.UserPhoneNumber, &u.UserRole); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByID retrieves a user by userId from the database.
// Returns a pointer to User and an error if not found or query fails.
func GetUserByID(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow(`SELECT userId, userFirstName, userLastName, userEmail, userPhoneNumber, userRole FROM users WHERE userId=?`, id).
		Scan(&u.UserID, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.UserPhoneNumber, &u.UserRole)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
