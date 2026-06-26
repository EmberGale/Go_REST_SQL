package config

import (
	"time"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Relay    RelayConfig    `mapstructure:"relay"`
	Redis    RedisConfig    `mapstructure:"redis"`
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

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	User         string        `mapstructure:"user"`
	DB           int           `mapstructure:"db"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	Timeout      time.Duration `mapstructure:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
