/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// LeaseSummary for dashboard leases section
type LeaseSummary struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Address        string `json:"address"`
	Unit           string `json:"unit"`
	LeaseStartDate string `json:"leaseStartDate"` // formatted date
}

func GetLeasesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT firstName, lastName, address, unit, leaseStartDate, rentAmount, leaseStatus FROM leasesView`)
		if err != nil {
			http.Error(w, "Failed to query leases", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var list []map[string]interface{}
		for rows.Next() {
			var firstName, lastName, address, unit, leaseStatus string
			var leaseStartDate int64
			var rentAmount int
			if err := rows.Scan(&firstName, &lastName, &address, &unit, &leaseStartDate, &rentAmount, &leaseStatus); err == nil {
				list = append(list, map[string]interface{}{
					"firstName":      firstName,
					"lastName":       lastName,
					"address":        address,
					"unit":           unit,
					"leaseStartDate": leaseStartDate,
					"rentAmount":     rentAmount,
					"leaseStatus":    leaseStatus,
				})
			}
		}
		json.NewEncoder(w).Encode(list)
	}
}
