package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {
	// Inject Note Handler
	authHandler := &AuthHandler{
		Engine: r,
		DB:     db,
	}
	authHandler.inject()

	user := &UserHandler{
		Engine: r,
		DB:     db,
	}

	user.inject()

}
