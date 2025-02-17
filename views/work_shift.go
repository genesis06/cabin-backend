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
	log.Debug("Get user id")
	log.Debug(userID)

	tx, err := database.DB.Begin()
	_, err = tx.Exec("INSERT INTO work_shifts(money_received, datetime, fk_user) VALUES ($1, $2, $3);", workShift.MoneyReceived, workShift.DateTime, userID)
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
	_, err = tx.Exec("UPDATE work_shifts SET money_received = $1, money_delivered= $2, datetime= $3, notes = $4 WHERE id = $5;", workShift.MoneyReceived, workShift.MoneyDelivered, workShift.DateTime, workShift.Notes, workShiftID)
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

func GetWorkShifts(c *gin.Context) {
	//limit := c.Query("limit")

	log.Println(c.Query("fromDate"))
	log.Println(c.Query("toDate"))

	sqlString := "SELECT ws.id, u.username, u.first_name, u.last_name, ws.money_received, ws.money_delivered, ws.datetime, ws.notes FROM work_shifts ws INNER JOIN users u ON u.id = ws.fk_user "

	if c.Query("fromDate") != "" && c.Query("toDate") != "" {
		sqlString = sqlString + "WHERE ws.datetime >= '" + c.Query("fromDate") + "' and ws.datetime <= '" + c.Query("toDate") + "' "
	}

	sqlString = sqlString + "OR ws.money_delivered=0 ORDER BY ws.id DESC"
	/*
		if limit != "" {
			sqlString += "LIMIT " + limit
		}
	*/
	log.Println(sqlString)

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	log.Debug("Get workShifts")

	workShifts := []*models.UserWorkShift{}
	for rows.Next() {
		var workShift models.UserWorkShift
		err := rows.Scan(&workShift.ID, &workShift.Username, &workShift.FirstName, &workShift.LastName, &workShift.MoneyReceived, &workShift.MoneyDelivered, &workShift.DateTime, &workShift.Notes)
		if err != nil {
			log.Fatal(err)
		}
		workShifts = append(workShifts, &workShift)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, workShifts)

}

func UpdateWorkShiftType(c *gin.Context) {

	var workShiftType models.WorkShiftType
	err := c.BindJSON(&workShiftType)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}
	log.Println(workShiftType)

	tx, err := database.DB.Begin()
	_, err = tx.Exec("UPDATE work_shifts_type SET start_time = $1, end_time = $2 WHERE id = $3;", workShiftType.StartTime, workShiftType.EndTime, workShiftType.ID)
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
