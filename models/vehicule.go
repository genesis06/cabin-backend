package models

type Vehicule struct {
	ID           int    `json:"id"`
	Type         int    `json:"type" binding:"required"`
	LicensePlate string `json:"license_plate"`
	Deleted      bool   `json:"deleted"`
}
