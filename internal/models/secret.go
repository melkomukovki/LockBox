package models

import "time"

type SecretType string

func (st SecretType) Validate() error {
	switch st {
	case Credentials, Text, Binary, PaymentCard, OTP:
		return nil
	default:
		return ErrSecretInvalidType
	}
}

// Secret types
const (
	Credentials SecretType = "credentials" // login and password
	Text        SecretType = "text"        // text data
	Binary      SecretType = "binary"      // any binary data
	PaymentCard SecretType = "card"        // card data. include: card number, date (month, year), cvv code
	OTP         SecretType = "otp"         // config string for otp
)

// Secret - secret model structure
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

type CredentialSecret struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PaymentCardSecret struct {
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
}

type TextSecret struct {
	Text string `json:"text"`
}

type BinarySecret struct {
	Binary string `json:"binary"`
}

type OTPSecret struct {
	OTP string `json:"otp"`
}
