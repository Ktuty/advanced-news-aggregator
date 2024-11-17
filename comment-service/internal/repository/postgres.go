package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// Подключение к БД
func NewPostgresDB(url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	return db, nil
}
