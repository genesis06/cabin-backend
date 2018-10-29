package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"log"

	"github.com/gin-gonic/gin"
)

func GetVehicules(c *gin.Context) {

	rentID := c.Param("id")

	log.Println(rentID)

	sqlString := "SELECT id, v_type, license_plate FROM vehicules WHERE fk_rent = " + rentID

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	vehicules := []*models.Vehicule{}
	for rows.Next() {
		var vehicule models.Vehicule
		err := rows.Scan(&vehicule.ID, &vehicule.Type, &vehicule.LicensePlate)
		if err != nil {
			log.Fatal(err)
		}
		vehicules = append(vehicules, &vehicule)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, vehicules)

}
