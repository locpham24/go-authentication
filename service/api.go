package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/handler"
)

type APIService struct {
	db *gorm.DB
}

func NewAPIService(db *gorm.DB) APIService {
	return APIService{
		db: db,
	}
}
func (s APIService) Start() {
	r := gin.Default()
	handler.InitRouter(r, s.db)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
