package config

import (
	"github.com/spf13/viper"
	"time"
)

// Flag vars
// TODO: move to main
//var (
//	configPath = flag.String("config", "config/config.yaml", "path to config file")
//)

type (
	Config struct {
		Server   ServerConfig   `mapstructure:"server"`
		Postgres PostgresConfig `mapstructure:"database"`
		Auth     AuthConfig     `mapstructure:"auth"`
		Log      LogConfig      `mapstructure:"log"`
	}

	ServerConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}

	AuthConfig struct {
		AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
		SigningKey     string        `mapstructure:"signing_key"`
	}

	LogConfig struct {
		Level string `mapstructure:"level"`
	}
)

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
