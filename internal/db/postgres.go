package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lerakravts/lite-calendar/internal/config"
	_ "github.com/lib/pq"
)

// Connect инициализирует соединение с PostgreSQL и возвращает *sqlx.DB
func Connect(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	// Настроим параметры пула соединений (best practices)
	db.SetMaxOpenConns(10)                  //не больше 10 соединений с БД
	db.SetMaxIdleConns(5)                   //5 соединений можно "держать наготове"
	db.SetConnMaxLifetime(30 * time.Minute) //Сколько живёт соединение, прежде чем его принудительно закрыть

	slog.Info("Connected to PostgreSQL", slog.String("host", cfg.PostgresHost), slog.String("db", cfg.PostgresDB))
	return db, nil
}
