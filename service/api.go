package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/handler"
	"github.com/locpham24/go-authentication/validator"
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
	binding.Validator = new(validator.DefaultValidator)
	handler.InitRouter(r, s.db)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
