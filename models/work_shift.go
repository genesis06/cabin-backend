package models

type WorkShift struct {
	ID             int    `json:"id"`
	MoneyReceived  int    `json:"money_received" binding:"required"`
	MoneyDelivered int    `json:"money_delivered"`
	DateTime       string `json:"date_time"`
	Username       string `json:"username" binding:"required"`
}
