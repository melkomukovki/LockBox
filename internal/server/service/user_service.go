package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/internal/server/validator"
)

var _ models.UserService = (*UserService)(nil)

type UserService struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) AuthUser(ctx context.Context, username, password string) (int, error) {
	if validator.ValidateUsername(username) != nil {
		return 0, models.ErrInvalidUsername
	}

	if validator.ValidatePassword(password) != nil {
		return 0, models.ErrInvalidPassword
	}

	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, models.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (u *UserService) RegisterUser(ctx context.Context, username, password string) (int, error) {
	if validator.ValidateUsername(username) != nil {
		return 0, models.ErrInvalidUsername
	}

	if validator.ValidatePassword(password) != nil {
		return 0, models.ErrInvalidPassword
	}

	_, err := u.repo.GetByUsername(ctx, username)
	if err == nil {
		return 0, models.ErrUserAlreadyExists
	} else if !errors.Is(err, models.ErrUserNotFound) {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userId, err := u.repo.Create(ctx, &models.User{Username: username, Password: string(hashedPassword)})
	if err != nil {
		return 0, err
	}

	return userId, nil
}
