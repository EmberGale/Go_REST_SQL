package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Relay    RelayConfig    `mapstructure:"relay"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type KafkaConfig struct {
	Acks     string   `mapstructure:"acks"`
	ClientId string   `mapstructure:"client_id"`
	Brokers  []string `mapstructure:"brokers"`
}

type RelayConfig struct {
	WorkerCount   int           `mapstructure:"worker_count"`
	BatchSize     int           `mapstructure:"batch_size"`
	PollInterval  time.Duration `mapstructure:"poll_interval"`
	MaxAttempts   int           `mapstructure:"max_attempts"`
	BaseDelay     int           `mapstructure:"base_delay"`
	MaxDelay      int           `mapstructure:"max_delay"`
	JitterPercent float64       `mapstructure:"jitter_percent"`
}

func Load() (*Config, error) {
	v := viper.New()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
		fmt.Println("Loaded configuration from .env file")
	} else {
		fmt.Println("No .env file found, using environment variables only")
	}

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	fmt.Printf("Viper env vars: %#v\n", v)

	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing .env file: %w", err)
	}
	/*
		if err := v.Unmarshal(&cfg, viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
			),
		)); err != nil {
			return nil, fmt.Errorf("unmarshaling config: %w", err)
		}
	*/
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("required env variable cfg.Database.Host is not set")
	}
	if cfg.Database.Port == 0 {
		return fmt.Errorf("required env variable cfg.Database.Port is not set")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("required env variable cfg.Database.User is not set")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("required env variable cfg.Database.Password is not set")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("required env variable cfg.Database.DBName is not set")
	}
	if cfg.Database.MigrationsPath == "" {
		return fmt.Errorf("required env variable cfg.Database.MigrationsPath is not set")
	}
	if cfg.Server.Port == "" {
		return fmt.Errorf("required env variable cfg.Server.Port is not set")
	}
	if cfg.Logger.Level == "" {
		return fmt.Errorf("required env variable cfg.Logger.Level is not set")
	}
	return nil
}
