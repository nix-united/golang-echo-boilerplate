package goose

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddEmailToUsers, downAddEmailToUsers)
}

func upAddEmailToUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE users
		ADD COLUMN email VARCHAR(200),
		ADD UNIQUE INDEX idx_users_email (email)
	`)
	return err
}

func downAddEmailToUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE users
		DROP COLUMN email
	`)
	return err
}
