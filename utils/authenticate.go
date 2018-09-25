package utils

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

// JWTSecret is the encryption string
var JWTSecret string

//InitiateTokenParams sets the JWT secret to be used to generate tokens
func InitiateTokenParams() {
	JWTSecret = viper.GetString("jwt.secret")
}

//GenerateTokens creates a new token for the authenticated user
func GenerateToken(claims map[string]interface{}) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		log.Println("error generando")
		return "", err
	}
	return tokenString, nil
}

//CheckJWTToken verifies if the provided token is valid
func CheckJWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(JWTSecret))
			return b, nil
		})
		if err == nil && token.Valid {
			c.Set("username", token.Claims["username"])
			c.Next()
			return
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.AbortWithError(400, errors.New("Invalid token"))
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				c.AbortWithError(400, errors.New("Expired token"))
				return
			} else {
				c.AbortWithError(400, ve)
				return
			}
		} else {
			c.AbortWithError(401, err)
			return
		}
	}
}
