package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	pb "github.com/locpham24/go-authentication/proto"
)

func InitRouter(r *gin.Engine, db *gorm.DB, client pb.AuthServiceClient) {
	// Inject Note Handler
	authHandler := &AuthHandler{
		Engine: r,
		DB:     db,
		client: client,
	}
	authHandler.inject()

	user := &UserHandler{
		Engine: r,
		DB:     db,
	}
	user.inject()

	// Inject health check handler
	healthCheckHandler := &HealthCheckHandler{
		Engine: r,
	}
	healthCheckHandler.inject()
}
