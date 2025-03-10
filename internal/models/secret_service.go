package models

import "context"

// ISecretService - интерфейс описывающий методы, необходимые для реализации на стороне сервиса
type ISecretService interface {
	CreateSecret(ctx context.Context, secret *Secret) (int, error)
	UpdateSecret(ctx context.Context, secret *Secret, userId int) error
	DeleteSecret(ctx context.Context, secretId, userId int) error
	GetSecret(ctx context.Context, secretId, userId int) (*Secret, error)
	GetAllSecrets(ctx context.Context, userId int) ([]Secret, error)
}
