package models

type VehiculeType struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}
