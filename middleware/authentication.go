package middleware

import (
	"github.com/locpham24/go-authentication/model"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. get token from header
		tokenString := c.Request.Header.Get("authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "user needs to be signed in to access this service",
			})
			c.Abort()
			return
		}

		// 2. Validate token
		token, _ := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if _, ok := token.Claims.(*model.Claims); ok && token.Valid {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			c.Abort()
		}

	}
}
