package models

import "time"

// User - модель описания пользователя
type User struct {
	ID        int       `gorm:"primary_key;auto_increment"`
	Username  string    `gorm:"not null;unique;size:32"`
	Password  string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
