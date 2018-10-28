package models

type WorkShift struct {
	ID             int    `json:"id"`
	MoneyReceived  int    `json:"money_received" binding:"required"`
	MoneyDelivered int    `json:"money_delivered"`
	DateTime       string `json:"date_time"`
	Username       string `json:"username" binding:"required"`
	Notes          string `json:"notes"`
}

type UserWorkShift struct {
	ID             int         `json:"id"`
	MoneyReceived  int         `json:"money_received"`
	MoneyDelivered int         `json:"money_delivered"`
	DateTime       interface{} `json:"date_time"`
	Username       string      `json:"username"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	Notes          interface{} `json:"notes"`
}
