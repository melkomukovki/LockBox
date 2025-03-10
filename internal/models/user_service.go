package models

import "context"

// IUserService - интерфейс описывающий необходимые методы для реализации сервиса
type IUserService interface {
	RegisterUser(ctx context.Context, username, password string) (int, error)
	AuthUser(ctx context.Context, username, password string) (int, error)
}
