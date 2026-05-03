package storage

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(storagePath string) error {
	m, err := migrate.New(
		"file://migrations",
		"sqlite3://"+storagePath,
	)
	if err != nil {
		log.Fatal("migration initialization error: ", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration up error: ", err)
	}
	return nil
}
