package models

import "context"

// UserService - интерфейс описывающий необходимые методы для реализации сервиса
type UserService interface {
	RegisterUser(ctx context.Context, username, password string) (int, error)
	AuthUser(ctx context.Context, username, password string) (int, error)
}
