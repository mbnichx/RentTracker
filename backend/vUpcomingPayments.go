/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file provides the dashboard upcoming payments handler, GetUpcomingRentHandler,
// which queries the upcomingPayments SQL view and returns upcoming rent payment data for the UI.

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type RentPayment struct {
	LeaseId         int     `json:"leaseId"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	Address         string  `json:"address"`
	Unit            string  `json:"unit"`
	RentAmount      float64 `json:"rentAmount"`
	LastPaymentUnix int64   `json:"lastPaymentUnix"`
	PaymentStatus   string  `json:"paymentStatus"` // "Overdue", "Due", "Paid"
}

// GetUpcomingRentHandler returns an HTTP handler for retrieving upcoming rent payment data from the upcomingPayments SQL view.
// Responds with a list of upcoming rent payments for the dashboard UI.
func GetUpcomingRentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT leaseId, firstName,lastName,address,unit,rentAmount,lastPaymentUnix,paymentStatus FROM upcomingPayments`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var payments []RentPayment
		for rows.Next() {
			var p RentPayment
			if err := rows.Scan(&p.LeaseId, &p.FirstName, &p.LastName, &p.Address, &p.Unit, &p.RentAmount, &p.LastPaymentUnix, &p.PaymentStatus); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			payments = append(payments, p)
		}
		json.NewEncoder(w).Encode(payments)
	}
}
