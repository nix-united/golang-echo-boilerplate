package list

import (
	mysql "github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

type UpdateUserTable struct{}

func (m *UpdateUserTable) GetName() string {
	return "UpdateUserTable"
}

func (m *UpdateUserTable) Up(con *sqlx.DB) {
	table := mysql.ChangeTable("users", con)
	table.String("email", 200).Unique()

	table.MustExec()
}

func (m *UpdateUserTable) Down(con *sqlx.DB) {
	table := mysql.ChangeTable("users", con)
	table.Column("email").Drop()

	table.MustExec()
}
