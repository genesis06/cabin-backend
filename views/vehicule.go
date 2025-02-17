package views

import (
	"cabin-backend/database"
	"cabin-backend/models"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func GetVehicules(c *gin.Context) {

	rentID := c.Param("id")

	log.Println(rentID)

	sqlString := "SELECT id, fk_vehicule_type, license_plate FROM rent_vehicules WHERE fk_rent = " + rentID

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	log.Debug("Get vehicules")

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

func GetVehiculeTypes(c *gin.Context) {
	sqlString := "SELECT id, name FROM vehicule_types;"

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	vehicules := []*models.VehiculeType{}
	for rows.Next() {
		var vehicule models.VehiculeType
		err := rows.Scan(&vehicule.ID, &vehicule.Name)
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
