package db

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func setup() error {
	goose.SetBaseFS(migrations)
	return goose.SetDialect("mysql")
}

// Up applies all pending migrations.
func Up(db *sql.DB) error {
	if err := setup(); err != nil {
		return err
	}
	return goose.Up(db, "migrations")
}

// Down rolls back the last applied migration.
func Down(db *sql.DB) error {
	if err := setup(); err != nil {
		return err
	}
	return goose.Down(db, "migrations")
}
