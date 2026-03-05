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
		return nil, fmt.Errorf("connection string is required")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("connection to database is failed %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping to database is failed %w", err)
	}

	slog.Info("database is connected")
	return db, nil
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		slog.Info("no connection to database")
	}
}
