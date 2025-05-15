package setup

import (
	"database/sql"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/migrations"

	"github.com/pressly/goose/v3"
)

func MigrateDB(db *sql.DB) error {
	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		return fmt.Errorf("set migrations dialect as mysql: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
