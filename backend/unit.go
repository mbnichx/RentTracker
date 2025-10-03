package backend

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// PROPERTY UNITS
type PropertyUnit struct {
	PropertyUnitID          int    `db:"propertyUnitId" json:"propertyUnitId"`
	PropertyID              int    `db:"propertyId" json:"propertyId"`
	PropertyUnitNumber      string `db:"propertyUnitNumber" json:"propertyUnitNumber"`
	PropertyUnitBeds        int    `db:"propertyUnitBeds" json:"propertyUnitBeds"`
	PropertyUnitBaths       int    `db:"propertyUnitBaths" json:"propertyUnitBaths"`
	PropertyUnitSqFt        int    `db:"propertyUnitSqFt" json:"propertyUnitSqFt"`
	PropertyUnitRentDefault int    `db:"propertyUnitRentDefault" json:"propertyUnitRentDefault"`
	PropertyUnitNotes       string `db:"propertyUnitNotes" json:"propertyUnitNotes"`
}

// == Handlers =====================================================================
// POST
func CreatePropertyUnitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var u PropertyUnit
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if u.PropertyID == 0 {
			respondError(w, http.StatusBadRequest, "propertyId required")
			return
		}
		id, err := CreatePropertyUnit(db, &u)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		u.PropertyUnitID = id
		respondJSON(w, http.StatusCreated, u)
	}
}

// GET
func GetPropertyUnitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/units/")
		if idStr == "" || idStr == "/" {
			units, err := GetAllPropertyUnits(db)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, units)
			return
		}
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		u, err := GetPropertyUnitByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not found")
			return
		}
		respondJSON(w, http.StatusOK, u)
	}
}

// PUT
func UpdatePropertyUnitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var u PropertyUnit
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if u.PropertyUnitID == 0 {
			respondError(w, http.StatusBadRequest, "propertyUnitId required")
			return
		}
		if err := UpdatePropertyUnit(db, &u); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DELETE
func DeletePropertyUnitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/units/")
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := DeletePropertyUnit(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries =====================================================================

func CreatePropertyUnit(db *sql.DB, u *PropertyUnit) (int, error) {
	res, err := db.Exec(`INSERT INTO propertyUnits (propertyId, propertyUnitNumber, propertyUnitBeds, propertyUnitBaths, propertyUnitSqFt, propertyUnitRentDefault, propertyUnitNotes)
	VALUES (?, ?, ?, ?, ?, ?, ?)`, u.PropertyID, u.PropertyUnitNumber, u.PropertyUnitBeds, u.PropertyUnitBaths, u.PropertyUnitSqFt, u.PropertyUnitRentDefault, u.PropertyUnitNotes)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func UpdatePropertyUnit(db *sql.DB, u *PropertyUnit) error {
	_, err := db.Exec(`UPDATE propertyUnits SET propertyId=?, propertyUnitNumber=?, propertyUnitBeds=?, propertyUnitBaths=?, propertyUnitSqFt=?, propertyUnitRentDefault=?, propertyUnitNotes=? WHERE propertyUnitId=?`,
		u.PropertyID, u.PropertyUnitNumber, u.PropertyUnitBeds, u.PropertyUnitBaths, u.PropertyUnitSqFt, u.PropertyUnitRentDefault, u.PropertyUnitNotes, u.PropertyUnitID)
	return err
}

func DeletePropertyUnit(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM propertyUnits WHERE propertyUnitId=?`, id)
	return err
}

func GetAllPropertyUnits(db *sql.DB) ([]PropertyUnit, error) {
	rows, err := db.Query(`SELECT propertyUnitId, propertyId, propertyUnitNumber, propertyUnitBeds, propertyUnitBaths, propertyUnitSqFt, propertyUnitRentDefault, propertyUnitNotes FROM propertyUnits`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PropertyUnit
	for rows.Next() {
		var u PropertyUnit
		if err := rows.Scan(&u.PropertyUnitID, &u.PropertyID, &u.PropertyUnitNumber, &u.PropertyUnitBeds, &u.PropertyUnitBaths, &u.PropertyUnitSqFt, &u.PropertyUnitRentDefault, &u.PropertyUnitNotes); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func GetPropertyUnitByID(db *sql.DB, id int) (*PropertyUnit, error) {
	var u PropertyUnit
	err := db.QueryRow(`SELECT propertyUnitId, propertyId, propertyUnitNumber, propertyUnitBeds, propertyUnitBaths, propertyUnitSqFt, propertyUnitRentDefault, propertyUnitNotes FROM propertyUnits WHERE propertyUnitId=?`, id).
		Scan(&u.PropertyUnitID, &u.PropertyID, &u.PropertyUnitNumber, &u.PropertyUnitBeds, &u.PropertyUnitBaths, &u.PropertyUnitSqFt, &u.PropertyUnitRentDefault, &u.PropertyUnitNotes)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
