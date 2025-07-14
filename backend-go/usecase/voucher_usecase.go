package usecase

import (
	"bookcabin-test/backend-go/model"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// ErrVoucherAlreadyExists indicates that a voucher already exists for the given flight and date.
// This error is returned when trying to generate vouchers for a flight that already has vouchers.
var ErrVoucherAlreadyExists = errors.New("vouchers already generated for this flight/date")

// Repository defines the interface that this usecase depends on.
// It abstracts the data layer operations for checking voucher existence and saving new vouchers.
type Repository interface {
	Exists(flightNumber, date string) (bool, error)
	Save(v model.Voucher) error
}

// VoucherUsecase holds the business logic for vouchers.
// It interacts with the repository to perform operations related to vouchers.
// This usecase is responsible for checking if a voucher exists and generating new vouchers based on business rules.
type VoucherUsecase struct {
	repo Repository
}

// NewVoucherUsecase creates a new VoucherUsecase.
// It takes a Repository as an argument, which is the data layer abstraction.
// This function is used to initialize the usecase with the repository, allowing it to perform business logic operations.
func NewVoucherUsecase(repo Repository) *VoucherUsecase {
	return &VoucherUsecase{repo: repo}
}

// CheckVoucher checks if a voucher exists for a specific flight and date.
// It uses the repository to query the database and returns true if a voucher exists, false otherwise.
func (uc *VoucherUsecase) CheckVoucher(flightNumber, date string) (bool, error) {
	return uc.repo.Exists(flightNumber, date)
}

// GenerateVouchers orchestrates the business logic for creating vouchers.
// It checks if a voucher already exists, generates unique seats, and saves the voucher to the repository.
func (uc *VoucherUsecase) GenerateVouchers(input model.Voucher) ([]string, error) {
	// Business Rule 1: Check if voucher already exists.
	exists, err := uc.repo.Exists(input.FlightNumber, input.FlightDate)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to check voucher existence: %w", err)
	}
	if exists {
		return nil, ErrVoucherAlreadyExists
	}

	// Business Rule 2: Generate 3 unique seats based on aircraft.
	seats, err := generateUniqueSeats(input.AircraftType)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to generate seats: %w", err)
	}
	input.Seats = seats

	// Business Rule 3: Save the voucher to the repository.
	err = uc.repo.Save(input)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to save voucher: %w", err)
	}

	return seats, nil
}

// It generates 3 unique seats based on the aircraft type.
// It uses randomization to ensure that the seats are unique and suitable for the specified aircraft.
func generateUniqueSeats(aircraftType string) ([]string, error) {
	// Initialize random seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var rows int
	var seatsPerRow []string

	switch aircraftType {
	case "ATR":
		rows = 18
		seatsPerRow = []string{"A", "C", "D", "F"}
	case "Airbus 320", "Boeing 737 Max":
		rows = 32
		seatsPerRow = []string{"A", "B", "C", "D", "E", "F"}
	default:
		return nil, fmt.Errorf("invalid aircraft type: %s", aircraftType)
	}

	// Use a map to track chosen seats for uniqueness
	chosenSeats := make(map[string]bool)
	var result []string

	// Generate 3 unique seats
	for len(result) < 3 {
		// Generate a random row and seat letter
		randomRow := r.Intn(rows) + 1 // +1 because rows are 1-based
		randomSeatIndex := r.Intn(len(seatsPerRow))
		seat := fmt.Sprintf("%d%s", randomRow, seatsPerRow[randomSeatIndex])

		// Check for uniqueness
		if !chosenSeats[seat] {
			chosenSeats[seat] = true
			result = append(result, seat)
		}
	}

	return result, nil
}
