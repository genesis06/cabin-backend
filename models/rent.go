package models

type Rent struct {
	ID               int        `json:"id"`
	CheckIn          string     `json:"check_in"`
	CheckOut         string     `json:"check_out"`
	CabinID          int        `json:"cabin_id" binding:"required"`
	ContratedTime    int        `json:"contracted_time" binding:"required"`
	Vehicules        []Vehicule `json:"vehicules" binding:"required"`
	Observations     string     `json:"observations"`
	NecessaryRepairs string     `json:"necessary_repairs"`
	LostStuff        string     `json:"lost_stuff"`
}

type LostStuff struct {
	Description string `json:"description" binding:"required"`
}
