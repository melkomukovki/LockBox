package models

import "context"

// IUserService - User service interface
type IUserService interface {
	RegisterUser(ctx context.Context, username, password string) (int, error)
	AuthUser(ctx context.Context, username, password string) (int, error)
}
