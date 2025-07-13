package store

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

// Ensure the SQLite driver is imported for database operations.
// DB is the global database connection.
var DB *sql.DB

// InitDB initializes the SQLite database connection and creates the vouchers table if it doesn't exist.
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	schema, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}

	if _, err := DB.Exec(string(schema)); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// CheckVoucherExists checks if a voucher already exists for the given flight number and date.
func CheckVoucherExists(flightNumber, date string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM vouchers WHERE flight_number = ? AND flight_date = ?)"
	err := DB.QueryRow(query, flightNumber, date).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SaveVoucher saves a new voucher with the provided details into the database.
// It expects three unique seats to be provided.
func SaveVoucher(name, id, flightNumber, date, aircraft string, seats []string) error {
	query := `INSERT INTO vouchers (crew_name, crew_id, flight_number, flight_date, aircraft_type, seat1, seat2, seat3)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, name, id, flightNumber, date, aircraft, seats[0], seats[1], seats[2])
	return err
}

// GenerateUniqueSeats generates three unique seats based on the aircraft type.
// It returns an error if the aircraft type is invalid.
func GenerateUniqueSeats(aircraftType string) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	var rows int
	var seatsPerRow []string

	// Define the number of rows and seat letters based on the aircraft type.
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

	// Generate all possible seats and shuffle them to get three unique random seats.
	if rows <= 0 || len(seatsPerRow) == 0 {
		return nil, fmt.Errorf("invalid aircraft configuration: %s", aircraftType)
	}
	var allSeats []string
	for r := 1; r <= rows; r++ {
		for _, s := range seatsPerRow {
			allSeats = append(allSeats, fmt.Sprintf("%d%s", r, s))
		}
	}

	// Ensure there are enough seats to select three unique ones.
	if len(allSeats) < 3 {
		return nil, fmt.Errorf("not enough seats available for aircraft type: %s", aircraftType)
	}
	// Shuffle the seats and select the first three.
	rand.Shuffle(len(allSeats), func(i, j int) {
		allSeats[i], allSeats[j] = allSeats[j], allSeats[i]
	})

	// Return the first three unique seats.
	if len(allSeats) < 3 {
		return nil, fmt.Errorf("not enough unique seats available for aircraft type: %s", aircraftType)
	}
	return allSeats[:3], nil
}
