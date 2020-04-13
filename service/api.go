package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"github.com/locpham24/go-authentication/handler"
	pb "github.com/locpham24/go-authentication/proto"
	"github.com/locpham24/go-authentication/validator"
	"github.com/micro/cli/v2"
	"google.golang.org/grpc"
)

type APIService struct {
	db *gorm.DB
}

func NewAPIService(db *gorm.DB) APIService {
	return APIService{
		db: db,
	}
}

func (s APIService) Start(ctx *cli.Context) {
	port := ctx.String("port")
	if port == "" {
		port = "7000"
	}

	fmt.Printf("Connect to server at :%s\n", port)
	// 1. Connect to server at TCP port
	conn, _ := grpc.Dial(fmt.Sprintf("localhost:%s", port), grpc.WithInsecure())
	// 2. New client
	client := pb.NewAuthServiceClient(conn)

	r := gin.Default()
	binding.Validator = new(validator.DefaultValidator)
	handler.InitRouter(r, s.db, client)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
