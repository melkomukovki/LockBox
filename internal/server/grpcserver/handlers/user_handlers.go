package handlers

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/pkg/auth"
)

// UserController описание структуры контроллера
type UserController struct {
	pb.UnimplementedUserServiceServer
	service    models.UserService
	jwtManager auth.JWTManager
}

// NewUserController конструктор для получения экземпляра контроллера
func NewUserController(service models.UserService, jwtManager auth.JWTManager) *UserController {
	return &UserController{service: service, jwtManager: jwtManager}
}

// SignUp - функция обработчик запросов на регистрацию пользователей
func (u *UserController) SignUp(ctx context.Context, request *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	_, err := u.service.RegisterUser(ctx, request.Login, request.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.SignUpResponse{Message: "success"}, nil
}

// SignIn - функция обработчик запросов на аутентификацию пользователей
func (u *UserController) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	userId, err := u.service.AuthUser(ctx, request.Login, request.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	accessToken, tokenTTL, err := u.jwtManager.NewJWT(strconv.Itoa(userId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SignInResponse{AccessToken: accessToken, ExpiresIn: int64(tokenTTL.Seconds())}, nil
}
