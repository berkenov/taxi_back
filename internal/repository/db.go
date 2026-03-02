package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// ConnectDB establishes connection to PostgreSQL
func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
