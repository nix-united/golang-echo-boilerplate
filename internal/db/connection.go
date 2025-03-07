package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
}

func (c Config) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

func NewConnection(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.dsn()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return db, nil
}
