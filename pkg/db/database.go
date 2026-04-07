package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"GoRestSQL/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// DB определяет интерфейс для работы с базой данных
type DB interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Beginx() (*sqlx.Tx, error)
	Close() error
	Ping() error
}

// Database реализует интерфейс DB
type Database struct {
	*sqlx.DB
}

// New создаёт новое подключение к БД, проверяет его и применяет миграции
func New(cfg config.DatabaseConfig) (*Database, error) {
	// Формируем DSN строку подключения
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Открываем подключение
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Применяем миграции
	migrationsPath := filepath.Clean(cfg.MigrationsPath)
	if err := goose.Up(db.DB, migrationsPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Database{db}, nil
}

// Close закрывает подключение к БД
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
