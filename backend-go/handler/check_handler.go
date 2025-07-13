package handler

import (
	"bookcabin-test/backend-go/store"
	"encoding/json"
	"log"
	"net/http"
)

// CheckRequest represents the request structure for checking voucher existence.
type CheckRequest struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
}

// CheckResponse represents the response structure for voucher existence check.
type CheckResponse struct {
	Exists bool `json:"exists"`
}

// CheckHandler handles the request to check if a voucher exists for a given flight number and date.
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := store.CheckVoucherExists(req.FlightNumber, req.Date)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("CHECK: Flight %s on %s - Exists: %t", req.FlightNumber, req.Date, exists)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CheckResponse{Exists: exists})
}
