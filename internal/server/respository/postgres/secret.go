package postgres

import (
	"context"
	"github.com/melkomukovki/LockBox/internal/models"
	"gorm.io/gorm"
)

var _ models.ISecretRepository = (*SecretRepository)(nil)

type SecretRepository struct {
	db *gorm.DB
}

func (s *SecretRepository) Create(ctx context.Context, secret *models.Secret) (int, error) {
	result := s.db.WithContext(ctx).Create(secret)
	return secret.ID, result.Error
}

func (s *SecretRepository) Delete(ctx context.Context, secretId int) error {
	result := s.db.WithContext(ctx).Where("id = ?", secretId).Delete(&models.Secret{})
	return result.Error
}

func (s *SecretRepository) Update(ctx context.Context, secret *models.Secret) error {
	result := s.db.WithContext(ctx).Model(&models.Secret{}).Where("id = ?", secret.ID).Updates(secret)
	return result.Error
}

func (s *SecretRepository) GetByUserId(ctx context.Context, userId int) ([]models.Secret, error) {
	var secrets []models.Secret
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&secrets)
	return secrets, result.Error
}

func (s *SecretRepository) Get(ctx context.Context, secretId int) (*models.Secret, error) {
	var secret models.Secret
	result := s.db.WithContext(ctx).Where("id = ?", secretId).First(&secret)
	return &secret, result.Error
}

func NewSecretRepository(db *gorm.DB) *SecretRepository {
	return &SecretRepository{db: db}
}
