package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// MaintenanceRequest matches your maintenanceRequests table
type MaintenanceRequestStatus struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Address           string `json:"address"`
	Unit              string `json:"unit"`
	Description       string `json:"description"`
	MaintenanceStatus string `json:"maintenanceStatus"`
	DateCreated       string `json:"dateCreated"` // formatted date
}

func GetMaintenanceRequestsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT firstName, lastName, address, unit, description, maintenanceStatus, dateCreated, priority, category FROM maintenanceRequestsView`)
		if err != nil {
			http.Error(w, "Failed to query maintenance requests", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var list []map[string]interface{}
		for rows.Next() {
			var firstName, lastName, address, unit, description, status, priority, category string
			var dateCreated int64
			if err := rows.Scan(&firstName, &lastName, &address, &unit, &description, &status, &dateCreated, &priority, &category); err == nil {
				list = append(list, map[string]interface{}{
					"firstName":   firstName,
					"lastName":    lastName,
					"address":     address,
					"unit":        unit,
					"description": description,
					"status":      status,
					"priority":    priority,
					"category":    category,
					"dateCreated": dateCreated,
				})
			}
		}
		json.NewEncoder(w).Encode(list)
	}
}
