/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file provides the dashboard overdue payments handler, GetOverduePaymentsHandler,
// which queries the overduePayments SQL view and returns overdue payment data for the UI.

import (
	"database/sql"
	"encoding/json"
	"github.com/guregu/null"
	"log"
	"net/http"
)

// OverduePayment represents a row from the overduePayments view
type OverduePayment struct {
	LeaseId         int      `json:"leaseId"`         // Lease identifier
	FirstName       string   `json:"firstName"`       // Tenant first name
	LastName        string   `json:"lastName"`        // Tenant last name
	Address         string   `json:"address"`         // Full property address (street, city, state, ZIP)
	Unit            string   `json:"unit"`            // Unit number
	RentAmount      float64  `json:"rentAmount"`      // Rent amount for this lease
	LastPaymentUnix null.Int `json:"lastPaymentUnix"` // Unix timestamp of last payment
	PaymentStatus   string   `json:"paymentStatus"`   // "Overdue", "Current", or "No payment recorded"
}

// GetOverduePaymentsHandler returns an HTTP handler for retrieving overdue payment data from the overduePayments SQL view.
// Responds with a list of overdue payments for the dashboard UI.
func GetOverduePaymentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT 
				leaseId, firstName, lastName, address, unit, rentAmount, 
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
				&o.LeaseId,
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
