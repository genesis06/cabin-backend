package views

import (
	"cabin-backend/database"
	"cabin-backend/models"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func GetReport(c *gin.Context) {

	log.Println(c.Query("fromDate"))
	log.Println(c.Query("toDate"))

	sqlString := "(SELECT CONCAT('Cabina ', c.cabin_number::text) as description, CONCAT(ct.quantity::text, ' horas') as amount, r.total as price, r.check_in::timestamp without time zone  as date_time FROM rents r INNER JOIN contracted_times ct ON ct.id = r.fk_contracted_time INNER JOIN cabins c ON c.id = r.fk_cabin "

	if c.Query("fromDate") != "" && c.Query("toDate") != "" {
		sqlString = sqlString + "WHERE r.check_in >= '" + c.Query("fromDate") + "' and r.check_in <= '" + c.Query("toDate") + "') "
	}

	sqlString = sqlString + "UNION (SELECT i.name, SUM(si.amount)::TEXT, SUM(si.price), null FROM sales s INNER JOIN sale_item si ON si.fk_sale = s.id INNER JOIN items i ON i.id = si.fk_item"

	if c.Query("fromDate") != "" && c.Query("toDate") != "" {
		sqlString = sqlString + " WHERE s.date > '" + c.Query("fromDate") + "' and s.date < '" + c.Query("toDate") + "'"
	}

	sqlString = sqlString + " GROUP BY i.name) ORDER BY date_time"

	log.Println(sqlString)

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	report := []*models.Report{}
	for rows.Next() {
		var item models.Report
		err := rows.Scan(&item.Description, &item.Amount, &item.Price, &item.DateTime)
		if err != nil {
			log.Fatal(err)
		}
		report = append(report, &item)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, report)
}

func GetSalesReport(c *gin.Context) {
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
