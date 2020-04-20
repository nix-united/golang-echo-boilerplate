package main

import (
	"echo-demo-project/migrations/list"
	"flag"
	gm "github.com/ShkrutDenis/go-migrations"
	gmStore "github.com/ShkrutDenis/go-migrations/store"
)

var isRollback *bool

func init() {
	isRollback = flag.Bool("rollback", false, "")
	flag.Parse()
}

func main() {
	if *isRollback {
		gm.Rollback(getMigrationsList())
		return
	}

	gm.Migrate(getMigrationsList())
}

func getMigrationsList() []gmStore.Migratable {
	return []gmStore.Migratable{
		&list.CreateUserTable{},
		&list.CreatePostTable{},
	}
}
