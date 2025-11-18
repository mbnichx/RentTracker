/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */

package main

// Package-level summary:
// This file contains common HTTP helpers for JSON response and error handling.
// Functions: respondJSON (writes JSON response), respondError (writes error as JSON).

import (
	"encoding/json"
	"net/http"
)

// Common helper for sending JSON responses
// respondJSON writes a JSON response with the given status code and payload.
// Sets Content-Type to application/json.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Common helper for error handling
// respondError writes an error message as a JSON response with the given status code.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
