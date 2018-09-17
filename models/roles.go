package models

type Role struct {
	Name string `json:"name" binding:"required"`
}
