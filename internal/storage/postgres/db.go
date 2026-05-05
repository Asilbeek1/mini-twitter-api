package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name)
	const op = "internal.db.New"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w ", op, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	return db, nil

}
