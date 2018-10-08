package models

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
}
