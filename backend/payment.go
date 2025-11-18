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
	"strconv"
	"strings"
)

// PAYMENTS
type Payment struct {
	PaymentID           int    `db:"paymentId" json:"paymentId"`
	LeaseID             int    `db:"leaseId" json:"leaseId"`
	PaymentAmount       int    `db:"paymentAmount" json:"paymentAmount"`
	PaymentDateUnix     int64  `db:"paymentDateUnix" json:"paymentDateUnix"`
	PaymentMethod       string `db:"paymentMethod" json:"paymentMethod"`
	PaymentNotes        string `db:"paymentNotes" json:"paymentNotes"`
	PaymentConfirmation []byte `db:"paymentConfirmation" json:"paymentConfirmation"`
}

// == Handlers =============================================================================
// POST
func CreatePaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var p Payment
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if p.LeaseID == 0 || p.PaymentAmount == 0 || p.PaymentDateUnix == 0 {
			respondError(w, http.StatusBadRequest, "leaseId, paymentAmount, paymentDateUnix required")
			return
		}
		id, err := CreatePayment(db, &p)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.PaymentID = id
		respondJSON(w, http.StatusCreated, p)
	}
}

// GET
func GetPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/payments/")
		if idStr == "" || idStr == "/" {
			list, err := GetAllPayments(db)
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
		p, err := GetPaymentByID(db, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not found")
			return
		}
		respondJSON(w, http.StatusOK, p)
	}
}

// PUT
func UpdatePaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var p Payment
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if p.PaymentID == 0 {
			respondError(w, http.StatusBadRequest, "paymentId required")
			return
		}
		if err := UpdatePayment(db, &p); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// DELETE
func DeletePaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/payments/")
		id, err := strconv.Atoi(strings.Trim(idStr, "/"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := DeletePayment(db, id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// == SQL Queries ======================================================================
func CreatePayment(db *sql.DB, p *Payment) (int, error) {
	res, err := db.Exec(`INSERT INTO payments (leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes, paymentConfirmation)
	VALUES (?, ?, ?, ?, ?, ?)`, p.LeaseID, p.PaymentAmount, p.PaymentDateUnix, p.PaymentMethod, p.PaymentNotes, p.PaymentConfirmation)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func UpdatePayment(db *sql.DB, p *Payment) error {
	_, err := db.Exec(`UPDATE payments SET leaseId=?, paymentAmount=?, paymentDateUnix=?, paymentMethod=?, paymentNotes=?, paymentConfirmation=? WHERE paymentId=?`,
		p.LeaseID, p.PaymentAmount, p.PaymentDateUnix, p.PaymentMethod, p.PaymentNotes, p.PaymentConfirmation, p.PaymentID)
	return err
}

func DeletePayment(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM payments WHERE paymentId=?`, id)
	return err
}

func GetAllPayments(db *sql.DB) ([]Payment, error) {
	rows, err := db.Query(`SELECT paymentId, leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes, paymentConfirmation FROM payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.PaymentID, &p.LeaseID, &p.PaymentAmount, &p.PaymentDateUnix, &p.PaymentMethod, &p.PaymentNotes, &p.PaymentConfirmation); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func GetPaymentByID(db *sql.DB, id int) (*Payment, error) {
	var p Payment
	err := db.QueryRow(`SELECT paymentId, leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes, paymentConfirmation FROM payments WHERE paymentId=?`, id).
		Scan(&p.PaymentID, &p.LeaseID, &p.PaymentAmount, &p.PaymentDateUnix, &p.PaymentMethod, &p.PaymentNotes, &p.PaymentConfirmation)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
