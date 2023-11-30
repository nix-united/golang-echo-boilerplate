package helpers

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm(dbMock *sql.DB) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: dbMock,
	}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}
