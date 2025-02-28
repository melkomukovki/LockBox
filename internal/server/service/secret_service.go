package service

import (
	"context"
	"github.com/melkomukovki/LockBox/internal/models"
)

var _ models.ISecretService = (*SecretService)(nil)

type SecretService struct {
	repo models.ISecretRepository
}

func (s *SecretService) CreateSecret(ctx context.Context, secret *models.Secret) (int, error) {
	if err := secret.Validate(); err != nil {
		return 0, err
	}

	secretId, err := s.repo.Create(ctx, secret)
	if err != nil {
		return 0, err
	}
	return secretId, nil
}

func (s *SecretService) UpdateSecret(ctx context.Context, secret *models.Secret, userId int) error {
	dbSecret, err := s.repo.Get(ctx, secret.ID)
	if err != nil {
		return err
	}

	if dbSecret.UserID != userId {
		return models.ErrSecretPermissionDenied
	}

	return s.repo.Update(ctx, secret)
}

func (s *SecretService) DeleteSecret(ctx context.Context, secretId, userId int) error {
	dbSecret, err := s.repo.Get(ctx, secretId)
	if err != nil {
		return err
	}

	if dbSecret.UserID != userId {
		return models.ErrSecretPermissionDenied
	}

	return s.repo.Delete(ctx, secretId)
}

func (s *SecretService) GetSecret(ctx context.Context, secretId, userId int) (*models.Secret, error) {
	secret, err := s.repo.Get(ctx, secretId)
	if err != nil {
		return nil, err
	}

	if secret.UserID != userId {
		return nil, models.ErrSecretPermissionDenied
	}

	return secret, nil
}

func (s *SecretService) GetAllSecrets(ctx context.Context, userId int) ([]models.Secret, error) {
	return s.repo.GetByUserId(ctx, userId)
}

func NewSecretService(repo models.ISecretRepository) *SecretService {
	return &SecretService{repo: repo}
}
