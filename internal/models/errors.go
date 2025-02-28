package models

import "errors"

var (
	// User errors
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("invalid password. 8 < len < 32")
	ErrInvalidUsername    = errors.New("invalid username. 5 < len < 32")

	// Secret errors
	ErrInvalidIdentifier      = errors.New("invalid identifier")
	ErrSecretAlreadyExists    = errors.New("secret already exists")
	ErrSecretNotFound         = errors.New("secret not found")
	ErrSecretPermissionDenied = errors.New("you do not have permission to access this secret")
	ErrSecretInvalidName      = errors.New("invalid secret name. length must be less than 32")
	ErrSecretInvalidType      = errors.New("invalid secret type")
	ErrSecretEmptyData        = errors.New("empty secret data")
)
