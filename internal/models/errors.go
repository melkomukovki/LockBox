package models

import "errors"

var (
	// Ошибки связанные с User
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("invalid password. 8 < len < 32")
	ErrInvalidUsername    = errors.New("invalid username. 5 < len < 32")

	// Ошибки связанные с Secret
	ErrSecretPermissionDenied = errors.New("you do not have permission to access this secret")
	ErrSecretInvalidName      = errors.New("invalid secret name. length must be less than 32")
	ErrSecretInvalidType      = errors.New("invalid secret type")
	ErrSecretEmptyData        = errors.New("empty secret data")
)
