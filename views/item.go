package views

import (
	"cabin-backend/database"
	"cabin-backend/models"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func GetSaleItems(c *gin.Context) {
	sqlString := "SELECT id, name, price FROM items ORDER BY id ASC"

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	items := []*models.Item{}
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, &item)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, items)
}
