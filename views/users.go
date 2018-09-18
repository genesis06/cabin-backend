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
	stmt, err := database.DB.Prepare("INSERT INTO users (username, first_name, last_name, password, status) VALUES ($1, $2, $3, $4, $5)")
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
