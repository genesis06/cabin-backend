package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"errors"
	"fmt"

	// "reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "github.com/davecgh/go-spew/spew"
)

// fix with tx, add roles
func CreateUser(c *gin.Context) {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, errors.New("Bad JSON"))
		return
	}
	fmt.Println(user)

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)
		return
	}

	tx, err := database.DB.Begin()
	stmt, err := tx.Prepare("INSERT INTO users (username, first_name, last_name, password, status) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		tx.Rollback()
		log.Error(err)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.AbortWithError(400, err)
		return
	}

	_, err = stmt.Exec(user.Username, user.FirstName, user.LastName, string(password), true)
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
	c.Data(201, gin.MIMEJSON, nil)
}

func GetUsers(c *gin.Context) {
	sqlString := "SELECT u.id, u.first_name, u.last_name, u.username, u.status, ws.start_time, ws.end_time FROM users u, work_shifts_type ws WHERE u.fk_work_shifts_type = ws.id ORDER BY id ASC"

	rows, err := database.DB.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Status, &user.StartTime, &user.EndTime)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(users); i++ {
		roles := []models.Role{}

		rows, err = database.DB.Query("SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.fk_role INNER JOIN users u ON ur.fk_user = u.id WHERE u.id = $1", users[i].ID)
		if err != nil {
			c.AbortWithError(500, err) //errors.New("Cant get rent"))
			return
		}

		for rows.Next() {
			role := models.Role{}
			err := rows.Scan(&role.ID, &role.Name)
			if err != nil {
				log.Fatal(err)
				c.AbortWithError(500, err)
				return
			}
			roles = append(roles, role)
		}
		users[i].Roles = roles
	}

	c.JSON(200, users)
}
