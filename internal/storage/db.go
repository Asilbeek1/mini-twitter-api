package storage

import (
	"database/sql"
	"fmt"
)

func New(storagePath string) (*sql.DB, error) {
	const op = "internal.db.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w ", op, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := runMigrations(storagePath); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	return &sql.DB{}, nil

}
