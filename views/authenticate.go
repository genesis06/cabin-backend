package views

import (
	"cabin-backend/database"
	"cabin-backend/models"
	"cabin-backend/utils"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type AuthenticationSchema struct {
	Username string
	Password string
}

func Authenticate(c *gin.Context) {
	var requestUser AuthenticationSchema
	var token map[string]string
	token = make(map[string]string)

	err := c.BindJSON(&requestUser)
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request"})
		return
	}

	err = utils.AuthenticateUser(requestUser.Username, requestUser.Password)
	if err != nil {
		log.Error(err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	if err != nil {
		log.Error(err)
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	log.Debug("Authenticating user")
	var user models.User

	err = database.DB.QueryRow("SELECT username, first_name, last_name, password FROM users WHERE username = $1", requestUser.Username).Scan(&user.Username, &user.FirstName, &user.LastName, &user.Password)
	if err != nil {
		log.Debug(err)
		c.AbortWithError(500, err)
		return
	}
	log.Debug("User info gotten")
	roles, err := database.DB.Query("SELECT r.name FROM users u,user_roles ur, roles r WHERE u.id = ur.fk_user AND r.id = ur.fk_role AND u.username = $1", user.Username)

	if err != nil {
		log.Fatal(err)
		c.AbortWithError(500, err)
		return
	}
	defer roles.Close()

	log.Debug("User roles gotten")

	rolesArray := []string{}

	for roles.Next() {
		role := ""
		err := roles.Scan(&role)
		if err != nil {
			log.Fatal(err)
			c.AbortWithError(500, err)
			return
		}
		rolesArray = append(rolesArray, role)
		log.Println(role)
	}
	err = roles.Err()
	if err != nil {
		log.Fatal(err)
	}

	claims := make(map[string]interface{})
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["username"] = user.Username
	claims["first_name"] = user.FirstName
	claims["last_name"] = user.LastName
	claims["roles"] = rolesArray
	tokenString, err := utils.GenerateToken(claims)

	log.Debug("Generating token")

	if err != nil {
		c.JSON(500, gin.H{"message": "Could not generate token"})
		log.Warning(err)
		return
	}

	token["token"] = tokenString
	c.JSON(200, token)
}
