package validator

import "github.com/melkomukovki/LockBox/internal/models"

// ValidatePassword функция для проверки паролей
// В дальнейшем можно добавить проверку на наличие букв с нижним/верхним регистром, спец символы и т.д.
func ValidatePassword(p string) error {
	if len(p) < 5 || len(p) > 32 {
		return models.ErrInvalidUsername
	}
	return nil
}
