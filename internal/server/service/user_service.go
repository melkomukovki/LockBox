package service

import (
	"context"
	"errors"
	"github.com/melkomukovki/LockBox/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var _ models.IUserService = (*UserService)(nil)

type UserService struct {
	repo models.IUserRepository
}

func NewUserService(repo models.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) AuthUser(ctx context.Context, username, password string) (int, error) {
	if err := validateCredentials(username, password); err != nil {
		return 0, err
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
	if err := validateCredentials(username, password); err != nil {
		return 0, err
	}

	_, err := u.repo.GetByUsername(ctx, username)
	if err == nil {
		return 0, models.ErrUserAlreadyExists
	} else if !errors.Is(err, models.ErrUserNotFound) {
		return 0, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	userId, err := u.repo.Create(ctx, &models.User{Username: username, Password: hashedPassword})
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func validateCredentials(username, password string) error {
	if len(username) < 5 || len(username) > 32 {
		return models.ErrInvalidUsername
	}
	if len(password) < 8 || len(password) > 32 {
		return models.ErrInvalidPassword
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
