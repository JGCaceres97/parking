package dtos

type EntryRequest struct {
	VehicleTypeID string `json:"vehicle_type_id"`
	LicensePlate  string `json:"license_plate"`
}

type ExitRequest struct {
	LicensePlate string `json:"license_plate"`
}
