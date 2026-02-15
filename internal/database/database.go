package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func (d *Database) GetDB() *sql.DB {
	return d.DB
}

func (d *Database) InitDB() error {
	connStr := os.Getenv("DATABASE_URL")

	var err error
	d.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	d.DB.SetMaxOpenConns(25)
	d.DB.SetMaxIdleConns(25)
	d.DB.SetConnMaxLifetime(5 * time.Minute)

	err = d.DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка ping БД: %v", err)
	}

	log.Println("База данных подключена успешно")
	return nil
}

func (d *Database) CloseDB() {
	if d.DB != nil {
		d.DB.Close()
		log.Println("Соединение с БД закрыто")
	}
}
