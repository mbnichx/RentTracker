package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type RentPayment struct {
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	Address         string  `json:"address"`
	Unit            string  `json:"unit"`
	RentAmount      float64 `json:"rentAmount"`
	LastPaymentUnix int64   `json:"lastPaymentUnix"`
	PaymentStatus   string  `json:"paymentStatus"` // "Overdue", "Due", "Paid"
}

func GetUpcomingRentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT firstName,lastName,address,unit,rentAmount,lastPaymentUnix,paymentStatus FROM upcomingPayments`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var payments []RentPayment
		for rows.Next() {
			var p RentPayment
			if err := rows.Scan(&p.FirstName, &p.LastName, &p.Address, &p.Unit, &p.RentAmount, &p.LastPaymentUnix, &p.PaymentStatus); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			payments = append(payments, p)
		}
		json.NewEncoder(w).Encode(payments)
	}
}
