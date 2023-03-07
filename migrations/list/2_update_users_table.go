package list

import (
	mysql "github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

const emailLength = 1000

type UpdateUserTable struct{}

func (m *UpdateUserTable) GetName() string {
	return "UpdateUserTable"
}

func (m *UpdateUserTable) Up(con *sqlx.DB) {
	table := mysql.ChangeTable("users", con)
	table.String("email", emailLength).Unique()

	table.MustExec()
}

func (m *UpdateUserTable) Down(con *sqlx.DB) {
	table := mysql.ChangeTable("users", con)
	table.Column("email").Drop()

	table.MustExec()
}
