package service

import (
	"context"
	"errors"

	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
)

// UserService структура User сервиса
type UserService struct {
	conn *grpcclient.Client
}

// Login функция для аутентификации пользователя
func (u *UserService) Login(ctx context.Context, login, password string) (string, error) {

	resp, err := u.conn.SignIn(ctx, login, password)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

// Register функция для регистрации пользователя
func (u *UserService) Register(ctx context.Context, login, password string) (string, error) {

	resp, err := u.conn.SignUp(ctx, login, password)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

// NewUserService конструктор для получения экземпляра User сервиса
func NewUserService(conn *grpcclient.Client) (*UserService, error) {
	if conn == nil {
		return nil, errors.New("grpc client must not be nil")
	}
	return &UserService{conn: conn}, nil
}
