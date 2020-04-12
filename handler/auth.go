package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/model"
	pb "github.com/locpham24/go-authentication/proto"
	"github.com/locpham24/go-authentication/validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
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

	req := pb.RegisterReq{
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    2000,
			"message": "Invalid password",
		})
		return
	}

	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	expirationTime := time.Now().Add(time.Duration(tokenTTL) * time.Minute)
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

func (a AuthHandler) refresh(c *gin.Context) {
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
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
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
	if !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    401,
			"message": "token is invalid",
		})
		return
	}

	if claims, ok := token.Claims.(*model.Claims); ok {
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 240*time.Second {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "too soon to refresh token",
			})
			return
		}

		tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

		expirationTime := time.Now().Add(time.Duration(tokenTTL) * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, tokenString)
		return
	}

	c.JSON(200, nil)
}
