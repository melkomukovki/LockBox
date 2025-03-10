package models

import "context"

// IUserRepository - интерфейс описывающий необходимые методы для реализации репозитория
type IUserRepository interface {
	Create(ctx context.Context, user *User) (int, error)
	GetById(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
