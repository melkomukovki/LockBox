package postgres

import (
	"context"
	"errors"
	"github.com/melkomukovki/LockBox/internal/models"
	"gorm.io/gorm"
)

var _ models.IUserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	tx := u.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, models.ErrUserAlreadyExists
		}
		return 0, err
	}

	return user.ID, tx.Commit().Error
}

func (u *UserRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
