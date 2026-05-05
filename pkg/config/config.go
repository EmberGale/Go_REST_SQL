package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

// DatabaseConfig содержит настройки подключения к БД
type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

// LoggerConfig содержит настройки логгера
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type KafkaConfig struct {
	BootstrapServers string `mapstructure:"bootstrapServers"`
	Acks             string `mapstructure:"acks"`
	ClientId         string `mapstructure:"clientId"`
}

// Load загружает конфигурацию из .env файла или переменных окружения
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	// Чтение переменных окружения с префиксом APP_
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Попытка прочитать .env файл (не критично, если не найдён)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Валидация обязательных полей
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig проверяет наличие обязательных переменных
func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("required env variable APP_DATABASE_HOST is not set")
	}
	if cfg.Database.Port == 0 {
		return fmt.Errorf("required env variable APP_DATABASE_PORT is not set")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("required env variable APP_DATABASE_USER is not set")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("required env variable APP_DATABASE_PASSWORD is not set")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("required env variable APP_DATABASE_DBNAME is not set")
	}
	if cfg.Database.MigrationsPath == "" {
		return fmt.Errorf("required env variable APP_DATABASE_MIGRATIONS_PATH is not set")
	}
	if cfg.Server.Port == "" {
		return fmt.Errorf("required env variable APP_SERVER_PORT is not set")
	}
	if cfg.Logger.Level == "" {
		return fmt.Errorf("required env variable APP_LOGGER_LEVEL is not set")
	}
	return nil
}
