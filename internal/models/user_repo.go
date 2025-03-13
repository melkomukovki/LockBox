package models

import "context"

// UserRepository - интерфейс описывающий необходимые методы для реализации репозитория
type UserRepository interface {
	Create(ctx context.Context, user *User) (int, error)
	GetById(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
