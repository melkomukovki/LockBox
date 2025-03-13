package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/melkomukovki/LockBox/internal/models"
)

const pgUniqueViolation = "23505"

var _ models.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, models.ErrUserAlreadyExists
		}
		return 0, err
	}

	return user.ID, nil
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
