package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/micro/cli/v2"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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

func (a *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.User, error) {
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

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AccessToken, error) {
	input := model.UserForm{
		Username: req.Username,
		Password: req.Password,
	}

	user := model.User{}
	user.Fill(input)

	// 2. verify if user exist
	err := a.DB.First(&user, "username = ?", user.Username).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, err
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
		return nil, err
	}

	resp := &pb.AccessToken{
		Token: tokenString,
	}

	return resp, nil
}

func (a *authService) Verify(ctx context.Context, req *pb.AccessToken) (*pb.User, error) {
	return nil, nil
}

func (a *authService) Refresh(ctx context.Context, req *pb.AccessToken) (*pb.AccessToken, error) {
	tokenString := req.Token
	if tokenString == "" {
		return nil, errors.New("user needs to be signed in to access this service")
	}

	// 2. Validate token
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	if claims, ok := token.Claims.(*model.Claims); ok {
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 240*time.Second {
			return nil, errors.New("too soon to refresh token")
		}

		tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

		expirationTime := time.Now().Add(time.Duration(tokenTTL) * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			return nil, err
		}

		res := &pb.AccessToken{Token: tokenString}
		return res, nil
	}

	return nil, errors.New("can not refresh token")
}

func (a authService) Start(ctx *cli.Context) {
	port := ctx.String("port")
	if port == "" {
		port = "7000"
	}
	// 1. Listen/Open a TPC connect at port
	listen, err := net.Listen("tcp", "ssl.local:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	certFile := "ssl.local.crt"
	keyFile := "ssl.local.key"

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	// 2. Tao server tu GRP
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 3. Map service to server
	pb.RegisterAuthServiceServer(grpcServer, &authService{
		DB: a.DB,
	})

	fmt.Printf("Starting gRPC server at ssl.local:%s .....", port)
	grpcServer.Serve(listen)
}
