package service

import (
	"context"
	"fmt"
	"log"

	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
)

// UserService структура User сервиса
type UserService struct {
	conn grpcclient.IGRPCClient
}

// Login функция для аутентификации пользователя
func (u UserService) Login(ctx context.Context) {
	fmt.Println("\tLogin `Page`!")
	login := inputString("Login > ")
	password := inputString("Password > ")

	resp, err := u.conn.SignIn(ctx, login, password)
	if err != nil {
		fmt.Printf("Failed to SignIn: %v\n", err)
		return
	}

	fmt.Printf("SignIn Success! Access Token: %s\n", resp.AccessToken)
}

// Register функция для регистрации пользователя
func (u UserService) Register(ctx context.Context) {
	fmt.Println("\tRegister `Page`!")
	login := inputString("Login > ")
	password := inputString("Password > ")
	confirmPassword := inputString("ConfirmPassword > ")

	if password != confirmPassword {
		fmt.Printf("Passwords do not match!\n")
		return
	}

	resp, err := u.conn.SignUp(ctx, login, password)
	if err != nil {
		fmt.Printf("Failed to SignUp: %v\n", err)
		return
	}

	fmt.Printf("SignUp Success! Response message: %s\n", resp.Message)
}

// NewUserService конструктор для получения экземпляра User сервиса
func NewUserService(conn grpcclient.IGRPCClient) IUserService {
	if conn == nil {
		log.Fatal("grpc  client must not be nil")
	}
	return &UserService{conn: conn}
}
