package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/middleware"
	"github.com/locpham24/go-authentication/model"
	"strconv"
)

type UserHandler struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func (u UserHandler) inject() {
	u.Engine.Use(middleware.AuthenticationRequired())
	u.Engine.GET("/user/:id", u.get)
}

func (u UserHandler) get(c *gin.Context) {
	idRaw := c.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	}
	user := model.User{}

	fmt.Println("hello")

	err = u.DB.First(&user, id).Error
	if err != nil {
		c.JSON(401, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, user)
}
