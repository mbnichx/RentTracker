package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// MAINTENANCE REQUESTS
type MaintenanceRequest struct {
	MaintenanceRequestID            int    `db:"maintenanceRequestId" json:"maintenanceRequestId"`
	PropertyUnitID                  int    `db:"propertyUnitId" json:"propertyUnitId"`
	LeaseID                         *int   `db:"leaseId" json:"leaseId,omitempty"`
	MaintenanceRequestInfo          string `db:"maintenanceRequestInfo" json:"maintenanceRequestInfo"`
	MaintenanceRequestPriority      string `db:"maintenanceRequestPriority" json:"maintenanceRequestPriority"`
	MaintenanceRequestCategory      string `db:"maintenanceRequestCategory" json:"maintenanceRequestCategory"`
	MaintenanceRequestStatus        string `db:"maintenanceRequestStatus" json:"maintenanceRequestStatus"`
	MaintenanceRequestCreatedUnix   int64  `db:"maintenanceRequestCreatedUnix" json:"maintenanceRequestCreatedUnix"`
	MaintenanceRequestCompletedUnix *int64 `db:"maintenanceRequestCompletedUnix" json:"maintenanceRequestCompletedUnix,omitempty"`
	MaintenanceRequestAssignedTo    string `db:"maintenanceAssignedTo" json:"maintenanceAssignedTo"`
}

// == Handlers ======================================================================================
func CreateMaintenanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var m MaintenanceRequest
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if m.PropertyUnitID == 0 || m.MaintenanceRequestInfo == "" || m.MaintenanceRequestCreatedUnix == 0 {
			respondError(w, http.StatusBadRequest, "propertyUnitId, info, createdUnix required")
			return
		}
		id, err := CreateMaintenanceRequest(db, &m)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		m.MaintenanceRequestID = id
		respondJSON(w, http.StatusCreated, m)
	}
}

func GetMaintenanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/maintenance/")
		if idStr == "" || idStr == "/" {
			list, err := GetAllMaintenanceRequests(db)
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
		m, err := GetMaintenanceRequestByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not found")
			return
		}
		respondJSON(w, http.StatusOK, m)
	}
}

func UpdateMaintenanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var m MaintenanceRequest
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if m.MaintenanceRequestID == 0 {
			respondError(w, http.StatusBadRequest, "maintenanceRequestId required")
			return
		}
		if err := UpdateMaintenanceRequest(db, &m); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func DeleteMaintenanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/maintenance/")
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := DeleteMaintenanceRequest(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries =====================================================================
func CreateMaintenanceRequest(db *sql.DB, m *MaintenanceRequest) (int, error) {
	res, err := db.Exec(`INSERT INTO maintenanceRequests (propertyUnitId, leaseId, maintenanceRequestInfo, maintenanceRequestPriority, maintenanceRequestCategory, maintenanceRequestStatus, maintenanceRequestCreatedUnix, maintenanceRequestCompletedUnix, maintenanceAssignedTo)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, m.PropertyUnitID, m.LeaseID, m.MaintenanceRequestInfo, m.MaintenanceRequestPriority, m.MaintenanceRequestCategory, m.MaintenanceRequestStatus, m.MaintenanceRequestCreatedUnix, m.MaintenanceRequestCompletedUnix, m.MaintenanceRequestAssignedTo)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}
func UpdateMaintenanceRequest(db *sql.DB, m *MaintenanceRequest) error {
	_, err := db.Exec(`UPDATE maintenanceRequests SET propertyUnitId=?, leaseId=?, maintenanceRequestInfo=?, maintenanceRequestPriority=?, maintenanceRequestCategory=?, maintenanceRequestStatus=?, maintenanceRequestCreatedUnix=?, maintenanceRequestCompletedUnix=?, maintenanceAssignedTo=? WHERE maintenanceRequestId=?`,
		m.PropertyUnitID, m.LeaseID, m.MaintenanceRequestInfo, m.MaintenanceRequestPriority, m.MaintenanceRequestCategory, m.MaintenanceRequestStatus, m.MaintenanceRequestCreatedUnix, m.MaintenanceRequestCompletedUnix, m.MaintenanceRequestAssignedTo, m.MaintenanceRequestID)
	return err
}

func DeleteMaintenanceRequest(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM maintenanceRequests WHERE maintenanceRequestId=?`, id)
	return err
}
func GetAllMaintenanceRequests(db *sql.DB) ([]MaintenanceRequest, error) {
	rows, err := db.Query(`SELECT maintenanceRequestId, propertyUnitId, leaseId, maintenanceRequestInfo, maintenanceRequestPriority, maintenanceRequestCategory, maintenanceRequestStatus, maintenanceRequestCreatedUnix, maintenanceRequestCompletedUnix, maintenanceAssignedTo FROM maintenanceRequests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MaintenanceRequest
	for rows.Next() {
		var m MaintenanceRequest
		if err := rows.Scan(&m.MaintenanceRequestID, &m.PropertyUnitID, &m.LeaseID, &m.MaintenanceRequestInfo, &m.MaintenanceRequestPriority, &m.MaintenanceRequestCategory, &m.MaintenanceRequestStatus, &m.MaintenanceRequestCreatedUnix, &m.MaintenanceRequestCompletedUnix, &m.MaintenanceRequestAssignedTo); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func GetMaintenanceRequestByID(db *sql.DB, id int) (*MaintenanceRequest, error) {
	var m MaintenanceRequest
	err := db.QueryRow(`SELECT maintenanceRequestId, propertyUnitId, leaseId, maintenanceRequestInfo, maintenanceRequestPriority, maintenanceRequestCategory, maintenanceRequestStatus, maintenanceRequestCreatedUnix, maintenanceRequestCompletedUnix, maintenanceAssignedTo FROM maintenanceRequests WHERE maintenanceRequestId=?`, id).
		Scan(&m.MaintenanceRequestID, &m.PropertyUnitID, &m.LeaseID, &m.MaintenanceRequestInfo, &m.MaintenanceRequestPriority, &m.MaintenanceRequestCategory, &m.MaintenanceRequestStatus, &m.MaintenanceRequestCreatedUnix, &m.MaintenanceRequestCompletedUnix, &m.MaintenanceRequestAssignedTo)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
