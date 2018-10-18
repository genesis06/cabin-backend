package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func CreateWorkShift(c *gin.Context) {
	var workShift models.WorkShift
	err := c.BindJSON(&workShift)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}
	log.Println(workShift)

	var userID int

	err = database.DB.QueryRow("SELECT id FROM users WHERE username = $1", workShift.Username).Scan(&userID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Cant get user"))
		return
	}
	log.Debug(userID)

	tx, err := database.DB.Begin()
	_, err = tx.Exec("INSERT INTO work_shifts(money_received, fk_user) VALUES ($1, $2);", workShift.MoneyReceived, userID)
	if err != nil {
		log.Println("ERRORRR")
		tx.Rollback()
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	c.Data(201, gin.MIMEJSON, nil)
}

func UpdateWorkShift(c *gin.Context) {
	workShiftID := c.Param("id")

	var workShift models.WorkShift
	err := c.BindJSON(&workShift)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}
	log.Println(workShift)

	tx, err := database.DB.Begin()
	_, err = tx.Exec("UPDATE work_shifts SET money_delivered= $1, datetime= $2 WHERE id = $3;", workShift.MoneyDelivered, workShift.DateTime, workShiftID)
	if err != nil {
		log.Println("ERRORRR")
		tx.Rollback()
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	c.Data(204, gin.MIMEJSON, nil)
}
