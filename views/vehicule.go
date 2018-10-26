package views

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetVehicules(c *gin.Context) {

	rentID := c.Param("id")

	log.Println(rentID)
	/*
		tx, err := database.DB.Begin()
		_ = tx.QueryRow("INSERT INTO sales(date) VALUES ($1) RETURNING id", sale.Date).Scan(&saleID)

		if err != nil {
			log.Println("ERRORRR 1")
			tx.Rollback()
			log.Println(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		for _, articule := range sale.SaleArticules {
			_, err = tx.Exec("INSERT INTO sale_item(fk_sale, fk_item, price, amount) VALUES ($1, $2, $3, $4)", saleID, articule.ArticuleID, articule.Price, articule.Amount)
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
	*/
}
