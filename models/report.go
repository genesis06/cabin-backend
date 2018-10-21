package models

type Report struct {
	Description string      `json:"description" binding:"required"`
	Amount      string      `json:"amount" binding:"required"`
	Price       int         `json:"price" binding:"required"`
	DateTime    interface{} `json:"date_time"`
}
