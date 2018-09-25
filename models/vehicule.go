package models

type Vehicule struct {
	ID           int    `json:"id"`
	Type         string `json:"type" binding:"required"`
	LicensePlate string `json:"license_plate"`
}
