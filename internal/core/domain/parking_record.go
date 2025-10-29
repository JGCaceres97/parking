package domain

import "time"

type ParkingRecord struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	VehicleTypeID   string     `json:"vehicle_type_id"`
	LicensePlate    string     `json:"license_plate"`
	EntryTime       time.Time  `json:"entry_time"`
	ExitTime        *time.Time `json:"exit_time"`
	TotalCharge     *float64   `json:"total_charge"`
	CalculatedHours *int       `json:"calculated_hours"`
}
