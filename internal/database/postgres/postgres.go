package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DB Struct to represent a PostgreSQL database connection
type DB struct {
	*sql.DB
}

// NewDB Function to create a new PostgreSQL database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	return &DB{db}, nil
}
