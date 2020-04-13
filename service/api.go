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
	"google.golang.org/grpc/credentials"
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
	certFile := "ssl.local.crt"

	fmt.Println("certFile:", certFile)
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Connect to server at ssl.local:%s\n", port)
	// 1. Connect to server at TCP port
	//conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", port), grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf("ssl.local:%s", port), grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err.Error())
	}

	// 2. New client
	client := pb.NewAuthServiceClient(conn)

	r := gin.Default()
	binding.Validator = new(validator.DefaultValidator)
	handler.InitRouter(r, s.db, client)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
