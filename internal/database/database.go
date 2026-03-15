package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	// Строка подключения (те же данные, что в docker-compose и goose)
	connStr := databaseURL
	// Открываем соединение
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	log.Println("Успешно подключено к PostgreSQL!")

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	return db, nil
}
