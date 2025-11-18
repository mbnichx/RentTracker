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
	"log"
	"net/http"
)

// OverduePayment represents a row from the overduePayments view
type OverduePayment struct {
	FirstName       string  `json:"firstName"`       // Tenant first name
	LastName        string  `json:"lastName"`        // Tenant last name
	Address         string  `json:"address"`         // Full property address (street, city, state, ZIP)
	Unit            string  `json:"unit"`            // Unit number
	RentAmount      float64 `json:"rentAmount"`      // Rent amount for this lease
	LastPaymentUnix int64   `json:"lastPaymentUnix"` // Unix timestamp of last payment
	PaymentStatus   string  `json:"paymentStatus"`   // "Overdue", "Current", or "No payment recorded"
}

func GetOverduePaymentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT 
				firstName, lastName, address, unit, rentAmount, 
				lastPaymentUnix, paymentStatus
			FROM overduePayments
		`)
		if err != nil {
			log.Printf("Error querying overduePayments: %v", err)
			http.Error(w, "Failed to retrieve overdue payments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var overdueList []OverduePayment
		for rows.Next() {
			var o OverduePayment
			if err := rows.Scan(
				&o.FirstName,
				&o.LastName,
				&o.Address,
				&o.Unit,
				&o.RentAmount,
				&o.LastPaymentUnix,
				&o.PaymentStatus,
			); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			overdueList = append(overdueList, o)
		}

		// Return JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(overdueList); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
