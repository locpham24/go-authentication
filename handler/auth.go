package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/model"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func (a AuthHandler) inject() {
	a.Engine.GET("/ping", a.ping)
	a.Engine.POST("/user/register", a.register)
	a.Engine.POST("/user/login", a.login)
}

func (a AuthHandler) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "pong",
	})
}

func (a AuthHandler) register(c *gin.Context) {
	input := model.UserForm{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"code":    1000,
			"message": err.Error(),
		})
		return
	}

	user := model.User{}
	user.Fill(input)

	// Check if user is existed
	var count int8
	err := a.DB.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error
	if err != nil {
		c.JSON(404, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
		return
	}

	// username is existed
	if count > 0 {
		c.JSON(400, gin.H{
			"code":    2000,
			"message": "username is existed",
		})
		return
	}

	// Add user to database
	err = a.DB.Create(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
		return
	}

	// generate token

	// response
	c.JSON(200, user)
}

func (a AuthHandler) login(c *gin.Context) {
	// 1. get user info from request
	input := model.UserForm{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"code":    1000,
			"message": err.Error(),
		})
		return
	}

	user := model.User{}
	user.Fill(input)

	// 2. verify if user exist
	err := a.DB.First(&user, "username = ?", user.Username).Error
	if err != nil {
		c.JSON(404, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
		return
	}

	// 3. check if password correct
	if input.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    2000,
			"message": "Invalid password",
		})
		return
	}

	// 4. create JWT token with expired time in 24h
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &model.Claims{
		Username: input.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1000,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, tokenString)
}
