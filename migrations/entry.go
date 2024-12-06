package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	_ "echo-demo-project/migrations/goose"
)

func main() {
	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	dir := flags.String("dir", "migrations/goose", "directory with migration files")

	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]

	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatal(err)
	}

	if err := goose.RunContext(context.Background(), command, db, *dir); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
