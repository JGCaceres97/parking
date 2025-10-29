package domain

type VehicleType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	HourlyRate  float64 `json:"hourly_rate"`
	Description string  `json:"description"`
}
