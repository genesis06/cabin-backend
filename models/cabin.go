package models

type Cabin struct {
	ID          int    `json:"id"`
	CabinNumber string `json:"cabin_number" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type CabinCheckout struct {
	ID          int        `json:"id binding:"required"`
	CabinNumber string     `json:"cabin_number" binding:"required"`
	CheckOut    string     `json:"check_out" binding:"required"`
	Vehicules   []Vehicule `json:"vehicules" binding:"required"`
}
