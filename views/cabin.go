package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Get cabins
func GetCabins(c *gin.Context) {
	sqlString := "SELECT c.id, c.cabin_number number, cs.name status FROM cabins c, cabin_status cs WHERE c.fk_status = cs.id ORDER BY c.id ASC"

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cabins := []*models.Cabin{}
	for rows.Next() {
		var cabin models.Cabin
		err := rows.Scan(&cabin.ID, &cabin.CabinNumber, &cabin.Status)
		if err != nil {
			log.Fatal(err)
		}
		cabins = append(cabins, &cabin)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, cabins)
}

func UpdateCabin(c *gin.Context) {
	var cabin models.Cabin
	cabinID := c.Param("id")

	err := c.BindJSON(&cabin)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}

	var statusID string
	err = database.DB.QueryRow("SELECT id FROM cabin_status WHERE name = $1", cabin.Status).Scan(&statusID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Cant get status"))
		return
	}

	tx, err := database.DB.Begin()
	stmt, err := database.DB.Prepare("UPDATE public.cabins SET fk_status= $1 WHERE id = $2;")
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)

		return
	}

	_, err = stmt.Exec(statusID, cabinID)
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

// Create rent
func CreateRent(c *gin.Context) {
	var rent models.Rent
	err := c.BindJSON(&rent)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}

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
	_ = tx.QueryRow("INSERT INTO rents( check_in, check_out, observations, necesary_repairs, fk_cabin, fk_contracted_time, total) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", rent.CheckIn, rent.CheckOut, rent.Observations, rent.NecessaryRepairs, rent.CabinID, contractedTimeID, price).Scan(&rentID)
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
