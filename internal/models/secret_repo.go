package models

import "context"

// ISecretRepository - secret repository interface
type ISecretRepository interface {
	Create(ctx context.Context, secret *Secret) (int, error)
	Delete(ctx context.Context, secretId int) error
	Update(ctx context.Context, secret *Secret) error
	GetByUserId(ctx context.Context, userId int) ([]Secret, error)
	Get(ctx context.Context, secretId int) (*Secret, error)
}
