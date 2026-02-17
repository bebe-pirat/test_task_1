package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(connStr string) (*sql.DB, error) {
	if connStr == "" {
		return nil, fmt.Errorf("Строка подключения не должна быть пустой")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Ошибка ping БД: %w", err)
	}

	slog.Info("База данных подключена успешно")
	return db, nil
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		slog.Info("Соединение с БД закрыто")
	}
}
