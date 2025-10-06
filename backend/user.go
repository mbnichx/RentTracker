package main

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func UpdateUser(db *sql.DB, u *User) error {
	_, err := db.Exec(`UPDATE users 
		SET userFirstName=?, userLastName=?, userEmail=?, userPhoneNumber=?, userRole=? 
		WHERE userId=?`,
		u.UserFirstName, u.UserLastName, u.UserEmail, u.UserPhoneNumber, u.UserRole, u.UserID)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE userId=?`, id)
	return err
}

// READ
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

func GetUserByID(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow(`SELECT userId, userFirstName, userLastName, userEmail, userPhoneNumber, userRole FROM users WHERE userId=?`, id).
		Scan(&u.UserID, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.UserPhoneNumber, &u.UserRole)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
