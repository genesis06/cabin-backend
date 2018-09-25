package models

type Cabin struct {
	ID          int    `json:"id"`
	CabinNumber string `json:"cabin_number" binding:"required"`
	Status      string `json:"status" binding:"required"`
}
