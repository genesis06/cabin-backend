package models

type Rent struct {
	ID               int        `json:"id"`
	CheckIn          string     `json:"check_in" binding:"required"`
	CheckOut         string     `json:"check_out" binding:"required"`
	CabinID          int        `json:"cabin_id" binding:"required"`
	ContratedTime    int        `json:"contracted_time" binding:"required"`
	Vehicules        []Vehicule `json:"vehicules" binding:"required"`
	Observations     string     `json:"observations"`
	NecessaryRepairs string     `json:"necessary_repairs"`
}
