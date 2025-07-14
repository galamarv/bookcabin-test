package repository

import (
	"bookcabin-test/backend-go/model"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "modernc.org/sqlite"
)

// VoucherRepository handles database operations for vouchers.
// It implements the Repository interface defined in the usecase layer.
// This repository is responsible for checking voucher existence and saving new vouchers to the database.
type VoucherRepository struct {
	db *sql.DB
}

// NewVoucherRepository creates a new repository instance.
// It takes a *sql.DB instance as an argument, which is the database connection.
// This function is used to initialize the repository with the database connection, allowing it to perform CRUD
func NewVoucherRepository(db *sql.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

// InitDB initializes the database connection and schema.
// It reads the schema from a file and creates the necessary tables if they do not exist.
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	schema, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("error reading schema file: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Database connection established and table verified.")
	return db, nil
}

// Exists checks if a voucher for a given flight and date exists.
// It queries the database to see if there is any record matching the flight number and date.
func (r *VoucherRepository) Exists(flightNumber, date string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM vouchers WHERE flight_number = ? AND flight_date = ?)"
	err := r.db.QueryRow(query, flightNumber, date).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Save saves a new voucher record to the database.
// It inserts the voucher details into the vouchers table.
// It expects the voucher to have all necessary fields populated, including crew name, ID, flight number, date, aircraft type, and seats.
// It returns an error if the insert operation fails.
func (r *VoucherRepository) Save(v model.Voucher) error {
	query := `INSERT INTO vouchers (crew_name, crew_id, flight_number, flight_date, aircraft_type, seat1, seat2, seat3)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, v.CrewName, v.CrewID, v.FlightNumber, v.FlightDate, v.AircraftType, v.Seats[0], v.Seats[1], v.Seats[2])
	return err
}
