package goose

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreatePosts, downCreatePosts)
}

func upCreatePosts(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE posts (
			id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(500) NULL,
			content VARCHAR(1000) NULL,
			deleted_at DATETIME NULL,
			user_id INT UNSIGNED,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
		)
	`)
	return err
}

func downCreatePosts(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS posts")
	return err
}
