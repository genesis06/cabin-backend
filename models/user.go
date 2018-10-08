package models

type User struct {
	ID        int    `db:"id" json:"id"`
	Username  string `db:"username" json:"username" binding:"required"`
	FirstName string `db:"first_name" json:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" binding:"required"`
	Password  string `db:"password" json:"password" binding:"required"`
	Roles     []Role `db:"roles" json:"roles" binding:"required"`
	Status    string `db:"status" json:"status"`
	StartTime string `db:"start_time" json:"start_time"`
	EndTime   string `db:"end_time" json:"end_time"`
}
