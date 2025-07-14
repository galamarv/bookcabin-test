package handler

import (
	"bookcabin-test/backend-go/model"
	"bookcabin-test/backend-go/usecase"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// VoucherUsecase defines the interface the handler depends on.
// This interface abstracts the usecase methods for checking and generating vouchers.
// It allows the handler to interact with the business logic without being tightly coupled to its implementation.
type VoucherUsecase interface {
	CheckVoucher(flightNumber, date string) (bool, error)
	GenerateVouchers(input model.Voucher) ([]string, error)
}

// Handler holds dependencies for the API handlers.
// It contains the usecase interface, which is injected during initialization.
// This allows the handler to call business logic methods without needing to know the details of their implementation
type Handler struct {
	usecase VoucherUsecase
}

// NewHandler creates a new handler with its dependencies.
// It takes a VoucherUsecase as an argument, which is the business logic layer.
// This function is used to initialize the handler with the usecase, allowing it to handle HTTP requests related to vouchers.
func NewHandler(uc VoucherUsecase) *Handler {
	return &Handler{usecase: uc}
}

// CheckVouchers handles the HTTP request for checking voucher existence.
func (h *Handler) CheckVouchers(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlightNumber string `json:"flightNumber"`
		Date         string `json:"date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := h.usecase.CheckVoucher(req.FlightNumber, req.Date)
	if err != nil {
		log.Printf("ERROR: Failed to check voucher: %v", err)
		http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
		return
	}

	log.Printf("CHECK: Flight %s on %s - Exists: %t", req.FlightNumber, req.Date, exists)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}

// GenerateVouchers handles the HTTP request for generating vouchers.
// It expects a JSON body with the crew name, ID, flight number, date, and aircraft type.
// It validates the input, checks for existing vouchers, generates new seats, and saves the voucher
func (h *Handler) GenerateVouchers(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		ID           string `json:"id"`
		FlightNumber string `json:"flightNumber"`
		Date         string `json:"date"`
		Aircraft     string `json:"aircraft"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	voucherInput := model.Voucher{
		CrewName:     req.Name,
		CrewID:       req.ID,
		FlightNumber: req.FlightNumber,
		FlightDate:   req.Date,
		AircraftType: req.Aircraft,
	}

	seats, err := h.usecase.GenerateVouchers(voucherInput)
	if err != nil {
		log.Printf("ERROR: Failed to generate vouchers: %v", err)
		if errors.Is(err, usecase.ErrVoucherAlreadyExists) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
		return
	}

	log.Printf("GENERATE: Assigned seats %v for flight %s", seats, req.FlightNumber)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "seats": seats})
}
