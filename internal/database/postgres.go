package database

import (
	"database/sql"
	"fmt"
	"horror_facts/internal/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgreDB cоздаёт подключение к БД (PostgreSQL)
func NewPostgreDB(cfg config.DataBaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("DB opening error: %w", err)
	}

	// Пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB error connecting: %w", err)
	}

	log.Println("Connected to PostgreSQL")

	return db, nil
}
