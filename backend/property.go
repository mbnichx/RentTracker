package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// PROPERTIES
type Property struct {
	PropertyID        int    `db:"propertyId" json:"propertyId"`
	OwnerUserID       int    `db:"ownerUserId" json:"ownerUserId"`
	PropertyName      string `db:"propertyName" json:"propertyName"`
	PropertyStreet    string `db:"propertyStreetAddress" json:"propertyStreetAddress"`
	PropertyCity      string `db:"propertyCity" json:"propertyCity"`
	PropertyState     string `db:"propertyState" json:"propertyState"`
	PropertyZip       string `db:"propertyZip" json:"propertyZip"`
	PropertyType      string `db:"propertyType" json:"propertyType"`
	PropertyYearBuilt int    `db:"propertyYearBuilt" json:"propertyYearBuilt"`
	PropertyNotes     string `db:"propertyNotes" json:"propertyNotes"`
}

// == Handlers ========================================================================
// POST
func CreatePropertyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var p Property
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		id, err := CreateProperty(db, &p)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.PropertyID = id
		respondJSON(w, http.StatusCreated, p)
	}
}

// GET
func GetPropertyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/properties/")
		if idStr == "" {
			props, err := GetAllProperties(db)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, props)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		prop, err := GetPropertyByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Not found")
			return
		}
		respondJSON(w, http.StatusOK, prop)
	}
}

// PUT
func UpdatePropertyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var p Property
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if p.PropertyID == 0 {
			respondError(w, http.StatusBadRequest, "Property ID is required")
			return
		}
		if err := UpdateProperty(db, &p); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DELETE
func DeletePropertyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/properties/delete/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		if err := DeleteProperty(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries ========================================================================
func CreateProperty(db *sql.DB, p *Property) (int, error) {
	res, err := db.Exec(`
		INSERT INTO properties (propertyName, propertyStreetAddress, propertyCity, propertyState, propertyZip, propertyOwnerId)
		VALUES (?, ?, ?, ?, ?, ?)`,
		p.PropertyName, p.PropertyStreet, p.PropertyCity, p.PropertyState, p.PropertyZip, p.OwnerUserID)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func UpdateProperty(db *sql.DB, p *Property) error {
	_, err := db.Exec(`UPDATE properties SET propertyName=?, propertyStreetAddress=?, propertyCity=?, propertyState=?, propertyZip=?, propertyOwnerId=? WHERE propertyId=?`,
		p.PropertyName, p.PropertyStreet, p.PropertyCity, p.PropertyState, p.PropertyZip, p.OwnerUserID, p.PropertyID)
	return err
}

func DeleteProperty(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM properties WHERE propertyId=?`, id)
	return err
}

func GetAllProperties(db *sql.DB) ([]Property, error) {
	rows, err := db.Query(`SELECT propertyId, propertyName, propertyAddress, propertyCity, propertyState, propertyZip, propertyOwnerId FROM properties`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var props []Property
	for rows.Next() {
		var p Property
		if err := rows.Scan(&p.PropertyID, &p.PropertyName, &p.PropertyStreet, &p.PropertyCity, &p.PropertyState, &p.PropertyZip, &p.OwnerUserID); err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	return props, nil
}

func GetPropertyByID(db *sql.DB, id int) (*Property, error) {
	var p Property
	err := db.QueryRow(`SELECT propertyId, propertyName, propertyStreetAddress, propertyCity, propertyState, propertyZip, propertyOwnerId FROM properties WHERE propertyId=?`, id).
		Scan(&p.PropertyID, &p.PropertyName, &p.PropertyStreet, &p.PropertyCity, &p.PropertyState, &p.PropertyZip, &p.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
