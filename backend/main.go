package backend

import (
	"database/sql"
	"log"
	"net/http"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./rt.db")
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	// Set pragmas for better defaults (optional)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal("enable fk:", err)
	}

	// // Initialize schema (creates tables if not exists)
	// if err := InitSchema(db); err != nil {
	// 	log.Fatal("init schema:", err)
	// }

	mux := http.NewServeMux()

	// helpers: POST endpoints for create, GET for listing/single,
	// PUT for update, DELETE for delete.
	// USERS
	mux.Handle("/users", CreateUserHandler(db))
	mux.Handle("/users/", GetUserHandler(db))
	mux.Handle("/users/update", UpdateUserHandler(db))
	mux.Handle("/users/delete/", DeleteUserHandler(db))

	// PROPERTIES
	mux.Handle("/properties", CreatePropertyHandler(db))
	mux.Handle("/properties/", GetPropertyHandler(db))
	mux.Handle("/properties/update", UpdatePropertyHandler(db))
	mux.Handle("/properties/delete/", DeletePropertyHandler(db))

	// UNITS
	mux.Handle("/units", CreatePropertyUnitHandler(db))
	mux.Handle("/units/", GetPropertyUnitHandler(db))
	mux.Handle("/units/update", UpdatePropertyUnitHandler(db))
	mux.Handle("/units/delete/", DeletePropertyUnitHandler(db))

	// TENANTS
	mux.Handle("/tenants", CreateTenantHandler(db))
	mux.Handle("/tenants/", GetTenantHandler(db))
	mux.Handle("/tenants/update", UpdateTenantHandler(db))
	mux.Handle("/tenants/delete/", DeleteTenantHandler(db))

	// LEASES
	mux.Handle("/leases", CreateLeaseHandler(db))
	mux.Handle("/leases/", GetLeaseHandler(db))
	mux.Handle("/leases/update", UpdateLeaseHandler(db))
	mux.Handle("/leases/delete/", DeleteLeaseHandler(db))

	// PAYMENTS
	mux.Handle("/payments", CreatePaymentHandler(db))
	mux.Handle("/payments/", GetPaymentHandler(db))
	mux.Handle("/payments/update", UpdatePaymentHandler(db))
	mux.Handle("/payments/delete/", DeletePaymentHandler(db))

	// MAINTENANCE
	mux.Handle("/maintenance", CreateMaintenanceHandler(db))
	mux.Handle("/maintenance/", GetMaintenanceHandler(db))
	mux.Handle("/maintenance/update", UpdateMaintenanceHandler(db))
	mux.Handle("/maintenance/delete/", DeleteMaintenanceHandler(db))

	// ACTIVITY LOGS
	mux.Handle("/activity", CreateActivityLogHandler(db))
	mux.Handle("/activity/", GetActivityLogHandler(db))
	mux.Handle("/activity/update", UpdateActivityLogHandler(db))
	mux.Handle("/activity/delete/", DeleteActivityLogHandler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
