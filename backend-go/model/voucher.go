package model

// Voucher represents the data structure for a voucher record.
type Voucher struct {
	ID           int64    `json:"id"`
	CrewName     string   `json:"crewName"`
	CrewID       string   `json:"crewId"`
	FlightNumber string   `json:"flightNumber"`
	FlightDate   string   `json:"flightDate"`
	AircraftType string   `json:"aircraftType"`
	Seats        []string `json:"seats"`
	CreatedAt    string   `json:"createdAt"` // ISO 8601 format
}
