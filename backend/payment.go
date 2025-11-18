/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file implements payment CRUD HTTP handlers and database helpers for
// creating, reading, updating, and deleting payment records. Handlers include
// CreatePaymentHandler, GetPaymentHandler, UpdatePaymentHandler, DeletePaymentHandler.
// DB helpers: CreatePayment, GetAllPayments, GetPaymentByID, UpdatePayment, DeletePayment.

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
// CreatePaymentHandler returns an HTTP handler for creating a new payment.
// Accepts a JSON body, validates required fields, inserts into DB, and responds with the created payment.
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
// GetPaymentHandler returns an HTTP handler for retrieving payments.
// If no ID is provided, returns all payments; otherwise, returns the payment with the given ID.
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
// UpdatePaymentHandler returns an HTTP handler for updating a payment.
// Accepts a JSON body, validates paymentId, updates the DB, and responds with status.
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
// DeletePaymentHandler returns an HTTP handler for deleting a payment by ID.
// Accepts a DELETE request, removes the payment from DB, and responds with status.
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
// CreatePayment inserts a new payment into the database.
// Returns the new payment ID and error if insertion fails.
func CreatePayment(db *sql.DB, p *Payment) (int, error) {
	res, err := db.Exec(`INSERT INTO payments (leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes, paymentConfirmation)
	VALUES (?, ?, ?, ?, ?, ?)`, p.LeaseID, p.PaymentAmount, p.PaymentDateUnix, p.PaymentMethod, p.PaymentNotes, p.PaymentConfirmation)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// UpdatePayment updates an existing payment in the database.
// Returns error if update fails.
func UpdatePayment(db *sql.DB, p *Payment) error {
	_, err := db.Exec(`UPDATE payments SET leaseId=?, paymentAmount=?, paymentDateUnix=?, paymentMethod=?, paymentNotes=?, paymentConfirmation=? WHERE paymentId=?`,
		p.LeaseID, p.PaymentAmount, p.PaymentDateUnix, p.PaymentMethod, p.PaymentNotes, p.PaymentConfirmation, p.PaymentID)
	return err
}

// DeletePayment removes a payment from the database by paymentId.
// Returns error if deletion fails.
func DeletePayment(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM payments WHERE paymentId=?`, id)
	return err
}

// GetAllPayments retrieves all payments from the database.
// Returns a slice of Payment and error if query fails.
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

// GetPaymentByID retrieves a payment by paymentId from the database.
// Returns pointer to Payment and error if not found or query fails.
func GetPaymentByID(db *sql.DB, id int) (*Payment, error) {
	var p Payment
	err := db.QueryRow(`SELECT paymentId, leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes, paymentConfirmation FROM payments WHERE paymentId=?`, id).
		Scan(&p.PaymentID, &p.LeaseID, &p.PaymentAmount, &p.PaymentDateUnix, &p.PaymentMethod, &p.PaymentNotes, &p.PaymentConfirmation)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
