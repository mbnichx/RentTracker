package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// TENANTS
type Tenant struct {
	TenantID        int    `db:"tenantId" json:"tenantId"`
	TenantFirstName string `db:"tenantFirstName" json:"tenantFirstName"`
	TenantLastName  string `db:"tenantLastName" json:"tenantLastName"`
	TenantEmail     string `db:"tenantEmailAddress" json:"tenantEmailAddress"`
	TenantPhone     string `db:"tenantPhoneNumber" json:"tenantPhoneNumber"`
}

// == Handlers =====================================================================
// POST
func CreateTenantHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode JSON body
		var t Tenant
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Basic validation
		if t.TenantFirstName == "" || t.TenantLastName == "" {
			respondError(w, http.StatusBadRequest, "Tenant first and last name required")
			return
		}
		if t.TenantEmail == "" {
			respondError(w, http.StatusBadRequest, "Tenant email required")
			return
		}

		// Insert into DB
		_, err := CreateTenant(db, &t)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, t)
	}
}

// GET
func GetTenantHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expect tenantId in query param
		tenantIDStr := r.URL.Query().Get("tenantId")
		if tenantIDStr == "" {
			respondError(w, http.StatusBadRequest, "Missing tenantId parameter")
			return
		}
		tenantID, err := strconv.Atoi(tenantIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid tenantId")
			return
		}

		tenant, err := GetTenantByID(db, tenantID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondError(w, http.StatusNotFound, "Tenant not found")
			} else {
				respondError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		respondJSON(w, http.StatusOK, tenant)
	}
}

// PUT
func UpdateTenantHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var t Tenant
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if t.TenantID == 0 {
			respondError(w, http.StatusBadRequest, "Tenant ID is required for update")
			return
		}

		err := UpdateTenant(db, &t)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DELETE
func DeleteTenantHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tenantIDStr := strings.TrimPrefix(r.URL.Path, "/tenants/")
		id, err := strconv.Atoi(tenantIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid tenantId")
			return
		}

		if err := DeleteTenant(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries =====================================================================

func CreateTenant(db *sql.DB, t *Tenant) (int, error) {
	res, err := db.Exec(`INSERT INTO tenants (tenantFirstName, tenantLastName, tenantEmailAddress, tenantPhoneNumber)
	VALUES (?, ?, ?, ?)`, t.TenantFirstName, t.TenantLastName, t.TenantEmail, t.TenantPhone)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}
func UpdateTenant(db *sql.DB, t *Tenant) error {
	_, err := db.Exec(`UPDATE tenants SET tenantFirstName=?, tenantLastName=?, tenantEmailAddress=?, tenantPhoneNumber=? WHERE tenantId=?`,
		t.TenantFirstName, t.TenantLastName, t.TenantEmail, t.TenantPhone, t.TenantID)
	return err
}

func DeleteTenant(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tenants WHERE tenantId=?`, id)
	return err
}

func GetAllTenants(db *sql.DB) ([]Tenant, error) {
	rows, err := db.Query(`SELECT tenantId, tenantFirstName, tenantLastName, tenantEmailAddress, tenantPhoneNumber FROM tenants`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Tenant
	for rows.Next() {
		var t Tenant
		if err := rows.Scan(&t.TenantID, &t.TenantFirstName, &t.TenantLastName, &t.TenantEmail, &t.TenantPhone); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func GetTenantByID(db *sql.DB, id int) (*Tenant, error) {
	var t Tenant
	err := db.QueryRow(`SELECT tenantId, tenantFirstName, tenantLastName, tenantEmailAddress, tenantPhoneNumber FROM tenants WHERE tenantId=?`, id).
		Scan(&t.TenantID, &t.TenantFirstName, &t.TenantLastName, &t.TenantEmail, &t.TenantPhone)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
