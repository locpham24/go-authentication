package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/model"
)

type AuthHandler struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func (a AuthHandler) inject() {
	a.Engine.POST("/user/register", a.register)
	a.Engine.POST("/user/login", a.login)
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
	c.JSON(200, nil)
}
