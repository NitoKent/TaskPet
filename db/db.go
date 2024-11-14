// db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int

	SSLMode string
}

func NewStorage(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not ping the database: %w", err)
	}

	log.Println("Connected to the database successfully.")
	return db, nil
}
