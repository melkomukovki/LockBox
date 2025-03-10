package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/internal/models/mocks"
)

func TestCreateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	svc := NewSecretService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		secret := &models.Secret{ID: 1, Name: "test name", Type: models.Credentials, UserID: 1, Data: []byte("test")}
		mockRepo.EXPECT().Create(gomock.Any(), secret).Return(1, nil)

		id, err := svc.CreateSecret(context.Background(), secret)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("ValidationError", func(t *testing.T) {
		secret := &models.Secret{}
		_, err := svc.CreateSecret(context.Background(), secret)
		assert.Error(t, err)
	})
}

func TestUpdateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	svc := NewSecretService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		secret := &models.Secret{ID: 1, UserID: 1, Data: []byte("updated")}
		mockRepo.EXPECT().Get(gomock.Any(), secret.ID).Return(secret, nil)
		mockRepo.EXPECT().Update(gomock.Any(), secret).Return(nil)

		err := svc.UpdateSecret(context.Background(), secret, 1)
		assert.NoError(t, err)
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		secret := &models.Secret{ID: 1, UserID: 1}
		mockRepo.EXPECT().Get(gomock.Any(), secret.ID).Return(secret, nil)

		err := svc.UpdateSecret(context.Background(), &models.Secret{ID: 1, UserID: 2}, 2)
		assert.ErrorIs(t, err, models.ErrSecretPermissionDenied)
	})
}

func TestDeleteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	svc := NewSecretService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().Get(gomock.Any(), 1).Return(&models.Secret{ID: 1, UserID: 1}, nil)
		mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)

		err := svc.DeleteSecret(context.Background(), 1, 1)
		assert.NoError(t, err)
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		mockRepo.EXPECT().Get(gomock.Any(), 1).Return(&models.Secret{ID: 1, UserID: 1}, nil)

		err := svc.DeleteSecret(context.Background(), 1, 2)
		assert.ErrorIs(t, err, models.ErrSecretPermissionDenied)
	})
}

func TestGetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	svc := NewSecretService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().Get(gomock.Any(), 1).Return(&models.Secret{ID: 1, UserID: 1}, nil)

		secret, err := svc.GetSecret(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, secret)
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		mockRepo.EXPECT().Get(gomock.Any(), 1).Return(&models.Secret{ID: 1, UserID: 1}, nil)

		_, err := svc.GetSecret(context.Background(), 1, 2)
		assert.ErrorIs(t, err, models.ErrSecretPermissionDenied)
	})
}

func TestGetAllSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	svc := NewSecretService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByUserId(gomock.Any(), 1).Return([]models.Secret{{ID: 1, UserID: 1}}, nil)

		secrets, err := svc.GetAllSecrets(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, secrets, 1)
	})
}
