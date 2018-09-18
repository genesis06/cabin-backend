package models

type User struct {
	Username  string `db:"username" json:"username" binding:"required"`
	FirstName string `db:"first_name" json:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" binding:"required"`
	Password  string `db:"password" json:"password" binding:"required"`
	Status    string `db:"status" json:"status"`
}
