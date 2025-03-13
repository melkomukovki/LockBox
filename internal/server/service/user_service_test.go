package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/internal/models/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIUserRepository(ctrl)
	service := NewUserService(repo)

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		setup     func()
		expect    int
		expectErr error
	}{
		{
			name: "Success",
			setup: func() {
				repo.EXPECT().GetByUsername(ctx, username).Return(&models.User{ID: 1, Username: username, Password: string(hashedPassword)}, nil)
			},
			expect:    1,
			expectErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				repo.EXPECT().GetByUsername(ctx, username).Return(nil, models.ErrUserNotFound)
			},
			expect:    0,
			expectErr: models.ErrInvalidCredentials,
		},
		{
			name: "Invalid password",
			setup: func() {
				repo.EXPECT().GetByUsername(ctx, username).Return(&models.User{ID: 1, Username: username, Password: "wronghash"}, nil)
			},
			expect:    0,
			expectErr: models.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			id, err := service.AuthUser(ctx, username, password)
			assert.Equal(t, tt.expect, id)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestUserService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIUserRepository(ctrl)
	service := NewUserService(repo)

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	tests := []struct {
		name      string
		setup     func()
		expect    int
		expectErr error
	}{
		{
			name: "Success",
			setup: func() {
				repo.EXPECT().GetByUsername(ctx, username).Return(nil, models.ErrUserNotFound)
				repo.EXPECT().Create(ctx, gomock.Any()).Return(1, nil)
			},
			expect:    1,
			expectErr: nil,
		},
		{
			name: "User already exists",
			setup: func() {
				repo.EXPECT().GetByUsername(ctx, username).Return(&models.User{ID: 1, Username: username}, nil)
			},
			expect:    0,
			expectErr: models.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			id, err := service.RegisterUser(ctx, username, password)
			assert.Equal(t, tt.expect, id)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
