CREATE TABLE IF NOT EXISTS vouchers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    crew_name TEXT NOT NULL,
    crew_id TEXT NOT NULL,
    flight_number TEXT NOT NULL,
    flight_date TEXT NOT NULL,
    aircraft_type TEXT NOT NULL,
    seat1 TEXT NOT NULL,
    seat2 TEXT NOT NULL,
    seat3 TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(flight_number, flight_date)
);