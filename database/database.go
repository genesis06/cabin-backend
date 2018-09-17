package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	// Postgres driver
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// DB is the database connection
var DB *sql.DB

// InitiateConnection initiates the connection to the database
func InitiateConnection() *sql.DB {
	var err error

	user := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	SSLMode := viper.GetString("database.sslmode")
	connectionString := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=" + SSLMode

	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Connection initialized")
	return DB
}

func getConnection() *sql.DB {
	return DB
}
