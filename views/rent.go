package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Create rent
func CreateRent(c *gin.Context) {
	var rent models.Rent
	err := c.BindJSON(&rent)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}
	log.Println(rent)
	var contractedTimeID string
	var price string

	err = database.DB.QueryRow("SELECT id, price FROM contracted_times WHERE quantity = $1", rent.ContratedTime).Scan(&contractedTimeID, &price)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Cant get status"))
		return
	}

	log.Println("timee: " + contractedTimeID)

	tx, err := database.DB.Begin()

	var rentID int
	_ = tx.QueryRow("INSERT INTO rents( check_in, observations, necesary_repairs, fk_cabin, fk_contracted_time, total) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", rent.CheckIn, rent.Observations, rent.NecessaryRepairs, rent.CabinID, contractedTimeID, price).Scan(&rentID)

	if err != nil {
		log.Println("ERRORRR 1")
		tx.Rollback()
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, vehicule := range rent.Vehicules {
		_, err = tx.Exec("INSERT INTO vehicules ( v_type, license_plate, fk_rent) VALUES ($1, $2, $3)", vehicule.Type, vehicule.LicensePlate, rentID)
		if err != nil {
			log.Println("ERRORRR 2")
			tx.Rollback()
			log.Println(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	tx.Commit()

	c.Data(201, gin.MIMEJSON, nil)
}

// GetRent get last rent of cabin
func GetRent(c *gin.Context) {
	cabinID := c.Param("id")

	var rent models.Rent

	err := database.DB.QueryRow("SELECT r.id, r.fk_cabin, r.check_in, r.observations, r.necesary_repairs, ct.quantity FROM rents r INNER JOIN contracted_times ct ON r.fk_contracted_time = ct.id INNER JOIN cabins c ON c.id = r.fk_cabin WHERE c.cabin_number = $1 ORDER BY r.id DESC LIMIT 1", cabinID).Scan(&rent.ID, &rent.CabinID, &rent.CheckIn, &rent.Observations, &rent.NecessaryRepairs, &rent.ContratedTime)
	if err != nil {
		c.AbortWithError(500, err) //errors.New("Cant get rent"))
		return
	}

	vehicules := []models.Vehicule{}

	rows, err := database.DB.Query("SELECT v.v_type, v.license_plate FROM vehicules v INNER JOIN rents r ON r.id = v.fk_rent INNER JOIN cabins c ON c.id = r.fk_cabin WHERE r.id = $1", rent.ID)
	if err != nil {
		c.AbortWithError(500, err) //errors.New("Cant get rent"))
		return
	}

	for rows.Next() {
		vehicule := models.Vehicule{}
		err := rows.Scan(&vehicule.Type, &vehicule.LicensePlate)
		if err != nil {
			log.Fatal(err)
			c.AbortWithError(500, err)
			return
		}
		vehicules = append(vehicules, vehicule)
		log.Println(vehicule)
	}

	rent.Vehicules = vehicules

	c.JSON(200, rent)
}

// Update observations and necessary repairs
func UpdateRent(c *gin.Context) {
	var rent models.Rent
	rentID := c.Param("id")

	err := c.BindJSON(&rent)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}

	tx, err := database.DB.Begin()
	stmt, err := database.DB.Prepare("UPDATE rents SET observations=$1, necesary_repairs=$2 WHERE id = $3;")
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)

		return
	}

	_, err = stmt.Exec(rent.Observations, rent.NecessaryRepairs, rentID)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)
		return
	}

	tx.Commit()
	//url := location.Get(c)
	//c.Header("Location", fmt.Sprintf("%s%s/%s", url, c.Request.URL, fmt.Sprintf("%d", lastID)))
	c.Data(204, gin.MIMEJSON, nil)
}

func PostCheckOut(c *gin.Context) {
	var cabinCheckout models.CabinCheckout
	cabinID := c.Param("id")

	err := c.BindJSON(&cabinCheckout)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}

	log.Println(cabinCheckout)

	var rentID int

	err = database.DB.QueryRow("SELECT id FROM rents WHERE fk_cabin = $1 ORDER BY id DESC LIMIT 1", cabinID).Scan(&rentID)
	if err != nil {
		c.AbortWithError(500, err) //errors.New("Cant get rent"))
		return
	}

	tx, err := database.DB.Begin()
	stmt, err := database.DB.Prepare("UPDATE rents SET check_out = $1 WHERE id = $2;")
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)

		return
	}

	_, err = stmt.Exec(cabinCheckout.CheckOut, rentID)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)
		return
	}

	tx.Commit()
	//url := location.Get(c)
	//c.Header("Location", fmt.Sprintf("%s%s/%s", url, c.Request.URL, fmt.Sprintf("%d", lastID)))
	c.Data(204, gin.MIMEJSON, nil)
}
