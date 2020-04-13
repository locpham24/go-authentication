package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/model"
	pb "github.com/locpham24/go-authentication/proto"
	"github.com/locpham24/go-authentication/validator"
	"net/http"
)

type AuthHandler struct {
	Engine *gin.Engine
	DB     *gorm.DB
	client pb.AuthServiceClient
}

func (a AuthHandler) inject() {
	a.Engine.POST("/register", a.register)
	a.Engine.POST("/login", a.login)
	a.Engine.GET("/refresh", a.refresh)
}

func (a AuthHandler) register(c *gin.Context) {
	input := model.UserForm{}
	if err := c.BindJSON(&input); err != nil {
		validator.HandleErrors(c, err)
		return
	}

	// Call to gRPC server
	req := pb.RegisterRequest{
		Username: input.Username,
		Password: input.Password,
	}
	res, err := a.client.Register(context.TODO(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
	}

	if res == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2000,
			"message": "Can not create",
		})
	}

	user := model.User{}
	// Check if user is existed
	err = a.DB.Model(&model.User{}).First(&user, res.Id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
		return
	}

	// response
	c.JSON(200, user)
}

func (a AuthHandler) login(c *gin.Context) {
	// 1. get user info from request
	input := model.UserForm{}
	if err := c.BindJSON(&input); err != nil {
		validator.HandleErrors(c, err)
		return
	}

	req := pb.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	}
	res, err := a.client.Login(context.TODO(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2000,
			"message": err.Error(),
		})
	}

	if res == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2000,
			"message": "Log in failed. Please try again",
		})
	}

	c.JSON(200, res.Token)
}

func (a AuthHandler) refresh(c *gin.Context) {
	tokenString := c.Request.Header.Get("authorization")
	req := &pb.AccessToken{
		Token: tokenString,
	}
	res, err := a.client.Refresh(context.TODO(), req)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res.Token)
}
