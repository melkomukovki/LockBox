// Package config содержит описание и методы для работы с конфигом сервиса
package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	// Config - общая структура конфига
	Config struct {
		Server   ServerConfig   `mapstructure:"server"`
		Postgres PostgresConfig `mapstructure:"database"`
		Auth     AuthConfig     `mapstructure:"auth"`
		Log      LogConfig      `mapstructure:"log"`
	}

	// ServerConfig параметры относящиеся к запуску сервера
	ServerConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	// PostgresConfig параметры необходимые для работы с базой данных
	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}

	// AuthConfig параметры необходимые для работы JWT менеджера
	AuthConfig struct {
		AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
		SigningKey     string        `mapstructure:"signing_key"`
	}

	// LogConfig параметры для конфигурации логгера
	LogConfig struct {
		Level string `mapstructure:"level"`
	}
)

// New функция для получения экземпляра конфига
func New(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
