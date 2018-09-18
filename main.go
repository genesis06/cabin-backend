package main

import (
	"app/utils"
	"cabin-backend/database"
	"cabin-backend/views"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tommy351/gin-cors"
)

// GetEngine returns the main engine
func GetEngine() *gin.Engine {
	router := gin.Default()
	router.Use(gin.ErrorLoggerT(gin.ErrorTypePrivate))
	// FIXME: no control or security
	router.Use(cors.Middleware(cors.Options{}))
	router.Use(location.Default())
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/authenticate", views.Authenticate)
	//	router.POST("/subscribe", views.Subscribe)
	v1 := router.Group("/api/cabin/v1", utils.CheckJWTToken())
	{
		//	v1.GET("/users", views.GetUsers)
		v1.POST("/users", views.CreateUser)
		//	v1.GET("/roles", views.GetRoles)
		//	v1.POST("/roles", views.CreateRole)
		//	v1.POST("/accounting", views.CreateLog)
		//	v1.GET("/items", views.GetItems)
		//	v1.POST("/items", views.CreateItem)
	}
	/*	users := v1.Group("/users")
		{
			users.GET("/:id", views.GetUser)
			users.DELETE("/:id", views.DeleteUser)
			users.PUT("/:id", views.UpdateUser)
			users.GET("/:id/roles", views.GetUserRoles)
			users.POST("/:id/roles", views.AssociateUserRoles)
			users.DELETE("/:id/roles/:role", views.DeleteUserRoles)
		}
		roles := v1.Group("/roles")
		{
			roles.GET("/:id", views.GetRole)
			roles.DELETE("/:id", views.DeleteRole)
			roles.PUT("/:id", views.UpdateRole)
			roles.GET("/:id/items", views.GetItemRoles)
			roles.POST("/:id/items", views.CreateItemRole)
			roles.DELETE("/:id/items/:item", views.DeleteItemRole)
		}
		items := v1.Group("/items")
		{
			items.GET("/:id", views.GetItem)
			items.PUT("/:id", views.UpdateItem)
			items.DELETE("/:id", views.DeleteItem)
		}

		recovery := v1.Group("/recovery")
		{
			recovery.PUT("/:id", views.ChangePassword)
			//items.PUT("/:id", views.RecoverPassword)
		}

		byname := v1.Group("/ByName")
		{
			byname.GET("roles/:name", views.GetRoleByName)
			byname.GET("roles/:name/items", views.GetItemRolesByName)
		}

		router.POST(viper.GetString("app.beurl")+"/recovery/password", views.ForgotPassword)
		router.POST(viper.GetString("app.beurl")+"/recovery/sign-up", views.SignUp)
	*/
	return router

}

func initializeApp() {
	viper.AddConfigPath("/opt") // Required for production deployment
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	viper.BindEnv("database.username", "PG_USERNAME")
	viper.BindEnv("database.password", "PG_PASSWORD")
	viper.BindEnv("database.host", "PG_PORT_5432_TCP_ADDR")
	viper.BindEnv("database.port", "PG_PORT_5432_TCP_PORT")
	viper.BindEnv("database.database", "PG_DB")
	viper.BindEnv("database.sslmode", "PG_SSL")

	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)

	utils.InitiateTokenParams()

}

func main() {
	initializeApp()
	db := database.InitiateConnection()
	defer db.Close()
	gin.SetMode(gin.ReleaseMode)
	router := GetEngine()
	router.Run(":8002") // listen and serve on 0.0.0.0:8080
}
