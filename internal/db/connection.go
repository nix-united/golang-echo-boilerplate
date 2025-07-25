package db

import (
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDB(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn(cfg)), &gorm.Config{
		Logger: newLoggerAdapter(),
	})
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return db, nil
}

func dsn(c config.DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}
