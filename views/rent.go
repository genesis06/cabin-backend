package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func GetRents(c *gin.Context) {
	limit := c.Query("limit")

	sqlString := "SELECT r.id, c.cabin_number, r.check_in, r.check_out, ct.quantity, r.observations, r.necesary_repairs, r.lost_stuff FROM rents r INNER JOIN contracted_times ct ON ct.id = r.fk_contracted_time INNER JOIN cabins c ON c.id = r.fk_cabin WHERE r.check_out IS NOT NULL ORDER BY check_in DESC "

	if limit != "" {
		sqlString += "LIMIT " + limit
	}

	log.Println(sqlString)

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rents := []*models.RentLostStuff{}
	for rows.Next() {
		var rent models.RentLostStuff
		err := rows.Scan(&rent.ID, &rent.CabinNumber, &rent.CheckIn, &rent.CheckOut, &rent.ContratedTime, &rent.Observations, &rent.NecessaryRepairs, &rent.LostStuff)
		if err != nil {
			log.Fatal(err)
		}
		rents = append(rents, &rent)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, rents)
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

	err := database.DB.QueryRow("SELECT r.id, r.fk_cabin, r.check_in, r.observations, r.necesary_repairs, ct.quantity, r.lost_stuff FROM rents r INNER JOIN contracted_times ct ON r.fk_contracted_time = ct.id INNER JOIN cabins c ON c.id = r.fk_cabin WHERE c.cabin_number = $1 ORDER BY r.id DESC LIMIT 1", cabinID).Scan(&rent.ID, &rent.CabinID, &rent.CheckIn, &rent.Observations, &rent.NecessaryRepairs, &rent.ContratedTime, &rent.LostStuff)
	if err != nil {
		c.AbortWithError(500, err) //errors.New("Cant get rent"))
		return
	}

	vehicules := []models.Vehicule{}

	rows, err := database.DB.Query("SELECT v.id, v.v_type, v.license_plate FROM vehicules v INNER JOIN rents r ON r.id = v.fk_rent INNER JOIN cabins c ON c.id = r.fk_cabin WHERE r.id = $1", rent.ID)
	if err != nil {
		c.AbortWithError(500, err) //errors.New("Cant get rent"))
		return
	}

	for rows.Next() {
		vehicule := models.Vehicule{}
		err := rows.Scan(&vehicule.ID, &vehicule.Type, &vehicule.LicensePlate)
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

	for _, vehicule := range rent.Vehicules {

		if vehicule.ID == 0 {
			log.Println("Insertó")
			_, err = tx.Exec("INSERT INTO vehicules ( v_type, license_plate, fk_rent) VALUES ($1, $2, $3)", vehicule.Type, vehicule.LicensePlate, rentID)
			if err != nil {
				log.Println("ERRORRR 2")
				tx.Rollback()
				log.Println(err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			log.Println("Actualizó")
			_, err = tx.Exec("UPDATE vehicules SET v_type = $1, license_plate = $2 WHERE id = $3;", vehicule.Type, vehicule.LicensePlate, vehicule.ID)
			if err != nil {
				log.Println("ERRORRR 2")
				tx.Rollback()
				log.Println(err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}

		/**/
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

func PostLostStuff(c *gin.Context) {
	var rent models.RentLostStuff

	rentID := c.Param("id")

	err := c.BindJSON(&rent)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Json"))
		return
	}

	log.Println(rent)

	tx, err := database.DB.Begin()
	stmt, err := database.DB.Prepare("UPDATE rents SET lost_stuff = $1, observations = $2, necesary_repairs = $3 WHERE id = $4;")
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)

		return
	}
	_, err = stmt.Exec(rent.LostStuff, rent.Observations, rent.NecessaryRepairs, rentID)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)
		return
	}

	for _, vehicule := range rent.Vehicules {

		if vehicule.ID == 0 {
			log.Println("Insertó")
			_, err = tx.Exec("INSERT INTO vehicules ( v_type, license_plate, fk_rent) VALUES ($1, $2, $3)", vehicule.Type, vehicule.LicensePlate, rentID)
			if err != nil {
				log.Println("ERRORRR 2")
				tx.Rollback()
				log.Println(err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			log.Println("Actualizó")
			_, err = tx.Exec("UPDATE vehicules SET v_type = $1, license_plate = $2 WHERE id = $3;", vehicule.Type, vehicule.LicensePlate, vehicule.ID)
			if err != nil {
				log.Println("ERRORRR 2")
				tx.Rollback()
				log.Println(err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}

		/**/
	}

	tx.Commit()
	c.Data(204, gin.MIMEJSON, nil)
}
