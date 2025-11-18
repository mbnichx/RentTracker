/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file implements the activity log HTTP handlers and database helpers
// used to create, read, update, and delete activity log records. Handlers
// exposed: CreateActivityLogHandler, GetActivityLogHandler,
// UpdateActivityLogHandler, DeleteActivityLogHandler. DB helpers include
// CreateActivityLog, GetAllActivityLogs, GetActivityLogByID, UpdateActivityLog,
// and DeleteActivityLog. The handlers validate input and use JSON request/response.

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ACTIVITY LOGS
type ActivityLog struct {
	LogID         int    `db:"logId" json:"logId"`
	UserID        int    `db:"userId" json:"userId"`
	EntityType    string `db:"entityType" json:"entityType"`
	EntityID      int    `db:"entityId" json:"entityId"`
	Action        string `db:"action" json:"action"`
	TimestampUnix int64  `db:"timestampUnix" json:"timestampUnix"`
}

// == Handlers ========================================================================
// POST
// CreateActivityLogHandler returns an HTTP handler for creating a new activity log entry.
// Accepts a JSON body, validates required fields, inserts into DB, and responds with the created log.
func CreateActivityLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var a ActivityLog
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if a.TimestampUnix == 0 {
			respondError(w, http.StatusBadRequest, "timestampUnix required")
			return
		}
		id, err := CreateActivityLog(db, &a)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		a.LogID = id
		respondJSON(w, http.StatusCreated, a)
	}
}

// GET
// GetActivityLogHandler returns an HTTP handler for retrieving activity logs.
// If no ID is provided, returns all logs; otherwise, returns the log with the given ID.
func GetActivityLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/activity/")
		if idStr == "" || idStr == "/" {
			list, err := GetAllActivityLogs(db)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, list)
			return
		}
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		a, err := GetActivityLogByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not found")
			return
		}
		respondJSON(w, http.StatusOK, a)
	}
}

// PUT
// UpdateActivityLogHandler returns an HTTP handler for updating an activity log entry.
// Accepts a JSON body, validates logId, updates the DB, and responds with status.
func UpdateActivityLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var a ActivityLog
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if a.LogID == 0 {
			respondError(w, http.StatusBadRequest, "logId required")
			return
		}
		if err := UpdateActivityLog(db, &a); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DELETE
// DeleteActivityLogHandler returns an HTTP handler for deleting an activity log entry by ID.
// Accepts a DELETE request, removes the log from DB, and responds with status.
func DeleteActivityLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/activity/")
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := DeleteActivityLog(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries ========================================================================
// CreateActivityLog inserts a new activity log into the database.
// Returns the new log ID and error if insertion fails.
func CreateActivityLog(db *sql.DB, a *ActivityLog) (int, error) {
	res, err := db.Exec(`INSERT INTO activityLogs (userId, entityType, entityId, action, timestampUnix) VALUES (?, ?, ?, ?, ?)`,
		a.UserID, a.EntityType, a.EntityID, a.Action, a.TimestampUnix)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// UpdateActivityLog updates an existing activity log in the database.
// Returns error if update fails.
func UpdateActivityLog(db *sql.DB, a *ActivityLog) error {
	_, err := db.Exec(`UPDATE activityLogs SET userId=?, entityType=?, entityId=?, action=?, timestampUnix=? WHERE logId=?`,
		a.UserID, a.EntityType, a.EntityID, a.Action, a.TimestampUnix, a.LogID)
	return err
}

// DeleteActivityLog removes an activity log from the database by logId.
// Returns error if deletion fails.
func DeleteActivityLog(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM activityLogs WHERE logId=?`, id)
	return err
}

// GetAllActivityLogs retrieves all activity logs from the database.
// Returns a slice of ActivityLog and error if query fails.
func GetAllActivityLogs(db *sql.DB) ([]ActivityLog, error) {
	rows, err := db.Query(`SELECT logId, userId, entityType, entityId, action, timestampUnix FROM activityLogs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActivityLog
	for rows.Next() {
		var a ActivityLog
		if err := rows.Scan(&a.LogID, &a.UserID, &a.EntityType, &a.EntityID, &a.Action, &a.TimestampUnix); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// GetActivityLogByID retrieves an activity log by logId from the database.
// Returns pointer to ActivityLog and error if not found or query fails.
func GetActivityLogByID(db *sql.DB, id int) (*ActivityLog, error) {
	var a ActivityLog
	err := db.QueryRow(`SELECT logId, userId, entityType, entityId, action, timestampUnix FROM activityLogs WHERE logId=?`, id).
		Scan(&a.LogID, &a.UserID, &a.EntityType, &a.EntityID, &a.Action, &a.TimestampUnix)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
