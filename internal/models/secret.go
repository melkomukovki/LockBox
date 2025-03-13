package models

import "time"

// SecretType тип секрета
type SecretType string

// Validate функция валидации типа секретов
func (st SecretType) Validate() error {
	switch st {
	case Credentials, Text, Binary, PaymentCard:
		return nil
	default:
		return ErrSecretInvalidType
	}
}

// Типы секретов
const (
	Credentials SecretType = "credentials" // Логин и пароль
	Text        SecretType = "text"        // Текстовая информация
	Binary      SecretType = "binary"      // Бинарная информация
	PaymentCard SecretType = "card"        // Данные платежной карты
)

// Secret - модель описывающая структуру секрета
type Secret struct {
	ID          int        `gorm:"primary_key;auto_increment"`
	Name        string     `gorm:"not null;size:32;index:,unique"`
	UserID      int        `gorm:"not null;index"`
	User        User       `gorm:"foreignKey:UserID"`
	Description string     `gorm:"size:512"`
	Type        SecretType `gorm:"not null"`
	Data        []byte     `gorm:"not null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
}

// Validate функция валидации данных секрета
func (s *Secret) Validate() error {
	if s.Name == "" || len(s.Name) > 32 {
		return ErrSecretInvalidName
	}
	if err := s.Type.Validate(); err != nil {
		return err
	}
	if len(s.Data) == 0 {
		return ErrSecretEmptyData
	}
	return nil
}

// CredentialSecret модель описывающая данные в секрете с типом Credentials
type CredentialSecret struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// PaymentCardSecret модель описывающая данные в секрете с типом PaymentCard
type PaymentCardSecret struct {
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
}

// TextSecret модель описывающая данные в секрете с типом Text
type TextSecret struct {
	Text string `json:"text"`
}

// BinarySecret модель описывающая данные в секрете с типом Binary
type BinarySecret struct {
	Binary string `json:"binary"`
}
