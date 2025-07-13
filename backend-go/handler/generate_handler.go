package handler

import (
	"bookcabin-test/backend-go/store"
	"encoding/json"
	"log"
	"net/http"
)

// GenerateRequest represents the request structure for generating a voucher.
type GenerateRequest struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	Aircraft     string `json:"aircraft"`
}

// GenerateResponse represents the response structure for voucher generation.
type GenerateResponse struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seats,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// GenerateHandler handles the request to generate a voucher for a flight.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := store.CheckVoucherExists(req.FlightNumber, req.Date)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(GenerateResponse{Success: false, Error: "Vouchers already generated for this flight/date."})
		return
	}

	seats, err := store.GenerateUniqueSeats(req.Aircraft)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(seats) < 3 {
		http.Error(w, "Not enough unique seats available", http.StatusInternalServerError)
		return
	}

	err = store.SaveVoucher(req.Name, req.ID, req.FlightNumber, req.Date, req.Aircraft, seats)
	if err != nil {
		http.Error(w, "Failed to save voucher", http.StatusInternalServerError)
		return
	}

	log.Printf("GENERATE: Assigned seats %v for flight %s", seats, req.FlightNumber)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(GenerateResponse{Success: true, Seats: seats})
}
