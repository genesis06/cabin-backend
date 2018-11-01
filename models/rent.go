package models

type Rent struct {
	ID                int         `json:"id"`
	CheckIn           string      `json:"check_in"`
	CheckOut          interface{} `json:"check_out"`
	EstimatedCheckout string      `json:"estimated_checkout"`
	CabinID           int         `json:"cabin_id" binding:"required"`
	ContratedTime     int         `json:"contracted_time" binding:"required"`
	Vehicules         []Vehicule  `json:"vehicules" binding:"required"`
	Observations      string      `json:"observations"`
	NecessaryRepairs  string      `json:"necessary_repairs"`
	LostStuff         interface{} `json:"lost_stuff"`
	SalesCheck        bool        `json:"sales_check"`
}

type NextCheckoutRents struct {
	CheckIn           string `json:"check_in"`
	EstimatedCheckout string `json:"estimated_checkout"`
	ContratedTime     int    `json:"contracted_time"`
	CabinNumber       string `json:"cabin_number"`
}

type LostStuff struct {
	Description string `json:"description" binding:"required"`
}

type RentLostStuff struct {
	ID                int         `json:"id"`
	CheckIn           string      `json:"check_in"`
	CheckOut          interface{} `json:"check_out"`
	EstimatedCheckout string      `json:"estimated_checkout"`
	Vehicules         []Vehicule  `json:"vehicules"`
	CabinNumber       string      `json:"cabin_number"`
	ContratedTime     int         `json:"contracted_time"`
	Observations      string      `json:"observations"`
	NecessaryRepairs  string      `json:"necessary_repairs"`
	LostStuff         interface{} `json:"lost_stuff"`
	SalesCheck        bool        `json:"sales_check"`
}
