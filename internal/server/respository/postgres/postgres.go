// Package postgres содержит код реализации уровня репозитория
// Для баз данных PostgreSQL
package postgres

import (
	"fmt"
	"github.com/melkomukovki/LockBox/internal/server/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
	"time"
)

// NewPgClient принимает параметры конфигурации для подключения к базе данных и возвращает соединение
func NewPgClient(cfg *config.PostgresConfig) (*gorm.DB, error) {
	dsn := (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.Database,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}).String()

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
