package models

import "context"

// ISecretRepository - интерфейс описывающий необходимые методы для реализации репозитория
type ISecretRepository interface {
	Create(ctx context.Context, secret *Secret) (int, error)
	Delete(ctx context.Context, secretId int) error
	Update(ctx context.Context, secret *Secret) error
	GetByUserId(ctx context.Context, userId int) ([]Secret, error)
	Get(ctx context.Context, secretId int) (*Secret, error)
}
