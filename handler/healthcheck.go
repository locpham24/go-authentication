package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckHandler struct {
	Engine *gin.Engine
}

func (h HealthCheckHandler) inject() {
	h.Engine.GET("/ping", h.ping)
}

func (h HealthCheckHandler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
