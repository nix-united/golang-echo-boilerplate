package helpers

import (
	"echo-demo-project/config"
	"echo-demo-project/server"
	"github.com/labstack/echo/v4"
)

func NewServer() *server.Server {
	s := &server.Server{
		Echo: echo.New(),
		DB:   Init(),
		Config: config.NewConfig(),
	}

	return s
}
