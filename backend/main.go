/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file is the main entry point for the RentTracker backend server. It
// opens the SQLite database, sets up HTTP routes for all API endpoints, and
// starts the server on port 8080. Route registration covers users, login,
// dashboard, rent, property, unit, tenant, lease, payment, maintenance, and activity log endpoints.

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

// main initializes the SQLite database, sets up HTTP routes for all API endpoints,
// and starts the RentTracker backend server on port 8080.
// It registers handlers for users, login, dashboard, rent, property, unit, tenant,
// lease, payment, maintenance, and activity log endpoints.
func main() {
	// Open SQLite database file
	db, err := sql.Open("sqlite", "../rt.db")
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	// Set pragmas for better defaults (optional)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal("enable fk:", err)
	}

	mux := http.NewServeMux()

	// Register new user endpoint
	mux.Handle("/users", CreateUserHandler(db))

	// Login endpoint
	mux.Handle("/login", LoginHandler(db))

	// Dashboard endpoints
	mux.Handle("/overduePayments", GetOverduePaymentsHandler(db))
	mux.Handle("/maintenanceRequestStatus", GetMaintenanceRequestsHandler(db))
	mux.Handle("/leaseOverview", GetLeasesHandler(db))

	// Rent endpoints
	mux.Handle("/upcomingPayments", GetUpcomingRentHandler(db))
	mux.Handle("/payments", CreatePaymentHandler(db))

	// User endpoints
	mux.Handle("/users/", GetUserHandler(db))
	mux.Handle("/users/me", GetCurrentUserHandler(db))
	mux.Handle("/users/update", UpdateUserHandler(db))
	mux.Handle("/users/delete/", DeleteUserHandler(db))

	// Property endpoints
	mux.Handle("/properties", CreatePropertyHandler(db))
	mux.Handle("/properties/", GetPropertyHandler(db))
	mux.Handle("/properties/update", UpdatePropertyHandler(db))
	mux.Handle("/properties/delete/", DeletePropertyHandler(db))

	// Unit endpoints
	mux.Handle("/units", CreatePropertyUnitHandler(db))
	mux.Handle("/units/", GetPropertyUnitHandler(db))
	mux.Handle("/units/update", UpdatePropertyUnitHandler(db))
	mux.Handle("/units/delete/", DeletePropertyUnitHandler(db))

	// Tenant endpoints
	mux.Handle("/tenants", CreateTenantHandler(db))
	mux.Handle("/tenants/", GetTenantHandler(db))
	mux.Handle("/tenants/update", UpdateTenantHandler(db))
	mux.Handle("/tenants/delete/", DeleteTenantHandler(db))

	// Lease endpoints
	mux.Handle("/leases", CreateLeaseHandler(db))
	mux.Handle("/leases/", GetLeaseHandler(db))
	mux.Handle("/leases/update", UpdateLeaseHandler(db))
	mux.Handle("/leases/delete/", DeleteLeaseHandler(db))

	// Payment endpoints
	mux.Handle("/payments/", GetPaymentHandler(db))
	mux.Handle("/payments/update", UpdatePaymentHandler(db))
	mux.Handle("/payments/delete/", DeletePaymentHandler(db))

	// Maintenance endpoints
	mux.Handle("/maintenance", CreateMaintenanceHandler(db))
	mux.Handle("/maintenance/", GetMaintenanceHandler(db))
	mux.Handle("/maintenance/update", UpdateMaintenanceHandler(db))
	mux.Handle("/maintenance/delete/", DeleteMaintenanceHandler(db))

	// Activity log endpoints
	mux.Handle("/activity", CreateActivityLogHandler(db))
	mux.Handle("/activity/", GetActivityLogHandler(db))
	mux.Handle("/activity/update", UpdateActivityLogHandler(db))
	mux.Handle("/activity/delete/", DeleteActivityLogHandler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
