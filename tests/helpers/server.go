package helpers

import (
	"echo-demo-project/internal/config"
	"echo-demo-project/internal/server"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *server.Server {
	s := &server.Server{
		Echo:   echo.New(),
		DB:     db,
		Config: config.NewConfig(),
	}

	return s
}
