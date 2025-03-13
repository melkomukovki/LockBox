package validator

import "github.com/melkomukovki/LockBox/internal/models"

// ValidateUsername функция проверки требований для имени пользователя
func ValidateUsername(u string) error {
	if len(u) < 5 || len(u) > 32 {
		return models.ErrInvalidUsername
	}
	return nil
}
