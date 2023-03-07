package list

import (
	mysql "github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

type CreateUserTable struct{}

func (m *CreateUserTable) GetName() string {
	return "CreateUserTable"
}

func (m *CreateUserTable) Up(con *sqlx.DB) {
	table := mysql.NewTable("users", con)
	table.Column("id").Type("int unsigned").Autoincrement()
	table.PrimaryKey("id")
	table.String("name", titleLength).Nullable()
	table.String("password", contentLength).Nullable()
	table.Column("deleted_at").Type("datetime").Nullable()
	table.WithTimestamps()

	table.MustExec()
}

func (m *CreateUserTable) Down(con *sqlx.DB) {
	mysql.DropTable("users", con).MustExec()
}
