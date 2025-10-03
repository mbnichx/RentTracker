package backend

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
func CreateActivityLog(db *sql.DB, a *ActivityLog) (int, error) {
	res, err := db.Exec(`INSERT INTO activityLogs (userId, entityType, entityId, action, timestampUnix) VALUES (?, ?, ?, ?, ?)`,
		a.UserID, a.EntityType, a.EntityID, a.Action, a.TimestampUnix)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func UpdateActivityLog(db *sql.DB, a *ActivityLog) error {
	_, err := db.Exec(`UPDATE activityLogs SET userId=?, entityType=?, entityId=?, action=?, timestampUnix=? WHERE logId=?`,
		a.UserID, a.EntityType, a.EntityID, a.Action, a.TimestampUnix, a.LogID)
	return err
}

func DeleteActivityLog(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM activityLogs WHERE logId=?`, id)
	return err
}

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

func GetActivityLogByID(db *sql.DB, id int) (*ActivityLog, error) {
	var a ActivityLog
	err := db.QueryRow(`SELECT logId, userId, entityType, entityId, action, timestampUnix FROM activityLogs WHERE logId=?`, id).
		Scan(&a.LogID, &a.UserID, &a.EntityType, &a.EntityID, &a.Action, &a.TimestampUnix)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
