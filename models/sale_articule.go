package models

type SaleArticule struct {
	ID         int `json:"id"`
	ArticuleID int `json:"articule_id"`
	Amount     int `json:"amount"`
	Price      int `json:"price"`
}
