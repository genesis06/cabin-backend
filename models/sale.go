package models

type Sale struct {
	ID            int            `json:"id"`
	Date          string         `json:"date" binding:"required"`
	SaleArticules []SaleArticule `json:"sale_articules" binding:"required"`
}
