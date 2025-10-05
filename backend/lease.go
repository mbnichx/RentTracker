package main

import (
	"database/sql"
	"encoding/json"

	// "errors"
	"net/http"
	"strconv"
	"strings"
)

// LEASES
type Lease struct {
	LeaseID              int    `db:"leaseId" json:"leaseId"`
	TenantID             int    `db:"tenantId" json:"tenantId"`
	PropertyUnitID       int    `db:"propertyUnitId" json:"propertyUnitId"`
	LeaseStartUnix       int64  `db:"leaseStartUnix" json:"leaseStartUnix"`
	LeaseEndUnix         *int64 `db:"leaseEndUnix" json:"leaseEndUnix,omitempty"`
	LeaseRentAmount      int    `db:"leaseRentAmount" json:"leaseRentAmount"`
	LeaseSecurityDeposit int    `db:"leaseSecurityDeposit" json:"leaseSecurityDeposit"`
	LeaseDocumentLink    string `db:"leaseDocumentLink" json:"leaseDocumentLink"`
	LeaseStatus          string `db:"leaseStatus" json:"leaseStatus"`
}

// == Handlers ====================================================================
// POST
func CreateLeaseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var l Lease
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if l.TenantID == 0 || l.PropertyUnitID == 0 || l.LeaseStartUnix == 0 || l.LeaseRentAmount == 0 {
			respondError(w, http.StatusBadRequest, "tenantId, propertyUnitId, leaseStartUnix, leaseRentAmount are required")
			return
		}
		id, err := CreateLease(db, &l)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		l.LeaseID = id
		respondJSON(w, http.StatusCreated, l)
	}
}

func GetLeaseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/leases/")
		if idStr == "" || idStr == "/" {
			list, err := GetAllLeases(db)
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
		l, err := GetLeaseByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not found")
			return
		}
		respondJSON(w, http.StatusOK, l)
	}
}

func UpdateLeaseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var l Lease
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if l.LeaseID == 0 {
			respondError(w, http.StatusBadRequest, "leaseId required")
			return
		}
		if err := UpdateLease(db, &l); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func DeleteLeaseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/leases/")
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := DeleteLease(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries =================================================================
func CreateLease(db *sql.DB, l *Lease) (int, error) {
	res, err := db.Exec(`INSERT INTO leases (tenantId, propertyUnitId, leaseStartUnix, leaseEndUnix, leaseRentAmount, leaseSecurityDeposit, leaseDocumentLink, leaseStatus)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, l.TenantID, l.PropertyUnitID, l.LeaseStartUnix, l.LeaseEndUnix, l.LeaseRentAmount, l.LeaseSecurityDeposit, l.LeaseDocumentLink, l.LeaseStatus)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func GetAllLeases(db *sql.DB) ([]Lease, error) {
	rows, err := db.Query(`SELECT leaseId, tenantId, propertyUnitId, leaseStartUnix, leaseEndUnix, leaseRentAmount, leaseSecurityDeposit, leaseDocumentLink, leaseStatus FROM leases`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Lease
	for rows.Next() {
		var l Lease
		if err := rows.Scan(&l.LeaseID, &l.TenantID, &l.PropertyUnitID, &l.LeaseStartUnix, &l.LeaseEndUnix, &l.LeaseRentAmount, &l.LeaseSecurityDeposit, &l.LeaseDocumentLink, &l.LeaseStatus); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func GetLeaseByID(db *sql.DB, id int) (*Lease, error) {
	var l Lease
	err := db.QueryRow(`SELECT leaseId, tenantId, propertyUnitId, leaseStartUnix, leaseEndUnix, leaseRentAmount, leaseSecurityDeposit, leaseDocumentLink, leaseStatus FROM leases WHERE leaseId=?`, id).
		Scan(&l.LeaseID, &l.TenantID, &l.PropertyUnitID, &l.LeaseStartUnix, &l.LeaseEndUnix, &l.LeaseRentAmount, &l.LeaseSecurityDeposit, &l.LeaseDocumentLink, &l.LeaseStatus)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func UpdateLease(db *sql.DB, l *Lease) error {
	_, err := db.Exec(`UPDATE leases SET tenantId=?, propertyUnitId=?, leaseStartUnix=?, leaseEndUnix=?, leaseRentAmount=?, leaseSecurityDeposit=?, leaseDocumentLink=?, leaseStatus=? WHERE leaseId=?`,
		l.TenantID, l.PropertyUnitID, l.LeaseStartUnix, l.LeaseEndUnix, l.LeaseRentAmount, l.LeaseSecurityDeposit, l.LeaseDocumentLink, l.LeaseStatus, l.LeaseID)
	return err
}

func DeleteLease(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM leases WHERE leaseId=?`, id)
	return err
}
