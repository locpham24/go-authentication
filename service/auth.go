package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/locpham24/go-authentication/model"
	pb "github.com/locpham24/go-authentication/proto"

	"google.golang.org/grpc"
)

type authService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) authService {
	return authService{
		DB: db,
	}
}

func (a *authService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.User, error) {
	input := model.UserForm{
		Username: req.Username,
		Password: req.Password,
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	input.Password = string(pass)

	user := model.User{}
	user.Fill(input)

	// Check if user is existed
	var count int8
	err = a.DB.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error
	if err != nil {
		return nil, err
	}

	// username is existed
	if count > 0 {
		return nil, errors.New("username is existed")
	}

	// Add user to database
	err = a.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	res := &pb.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}
	return res, nil
}

func (a *authService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	return nil, nil
}

func (a *authService) Verify(ctx context.Context, req *pb.AccessToken) (*pb.User, error) {
	return nil, nil
}

func (s authService) Start() {
	// 1. Listen/Open a TPC connect at port
	lis, _ := net.Listen("tcp", ":50051")
	// 2. Tao server tu GRP
	grpcServer := grpc.NewServer()
	// 3. Map service to server
	pb.RegisterAuthServiceServer(grpcServer, &authService{
		DB: s.DB,
	})
	// 4. Binding port
	reflection.Register(grpcServer)
	fmt.Println("Start service")
	grpcServer.Serve(lis)
}
