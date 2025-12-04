package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

// Config хранит конфигурацию для подключения к базе данных.
type Config struct {
	DSN             string        // Строка подключения (Data Source Name)
	MaxOpenConns    int           // Максимальное количество открытых соединений
	MaxIdleConns    int           // Максимальное количество простаивающих соединений
	ConnMaxIdleTime time.Duration // Максимальное время жизни простаивающего соединения
}

// NewConnection инициализирует и возвращает пул соединений с PostgreSQL.
func NewConnection(cfg Config) (*sql.DB, error) {
	// Устанавливаем значения по умолчанию, если они не заданы
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 25
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = 5 * time.Minute
	}

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Проверяем соединение с базой данных с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// Migrate выполняет миграции схемы базы данных.
// В реальном проекте здесь лучше использовать специализированные инструменты,
// например, golang-migrate/migrate.
func Migrate(ctx context.Context, db *sql.DB) error {
	// Транзакция гарантирует, что все DDL-запросы выполнятся успешно, либо ни один.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	// Гарантируем откат транзакции в случае ошибки.
	defer tx.Rollback()

	schema := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = tx.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("error executing migration schema: %w", err)
	}

	// Если все прошло успешно, подтверждаем транзакцию.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	slog.Info("Database migration completed successfully")
	return nil
}
