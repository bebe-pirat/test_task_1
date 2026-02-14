package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка ping БД: %v", err)
	}

	log.Println("База данных подключена успешно")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Соединение с БД закрыто")
	}
}
