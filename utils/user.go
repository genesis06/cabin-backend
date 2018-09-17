package utils

import (
	"cabin-backend/database"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// AuthenticateUser is ...
func AuthenticateUser(username string, password string) error {
	//var user models.User

	var status bool
	var passwordHash string

	//db := database.InitiateConnection()
	err := database.DB.QueryRow("SELECT status FROM users WHERE username = $1", username).Scan(&status)

	if err != nil {
		return err
	}

	if !status {
		err := fmt.Errorf(`User "%v" is inactive`, username)
		log.Debug(err)
		return err
	}

	err = database.DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&passwordHash)
	if err != nil {
		return err
	}

	if password != passwordHash {
		err := fmt.Errorf(`Incorrect password`)
		log.Debug(err)
		return err
	}

	return nil
}
