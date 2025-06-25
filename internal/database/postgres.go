package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB(addr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, fmt.Errorf("database connect failed: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("database connected")
	return db, nil
}
