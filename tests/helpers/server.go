package helpers

import (
	"echo-demo-project/server"
	"github.com/labstack/echo/v4"
)

func NewServer() *server.Server {
	s := &server.Server{
		Echo: echo.New(),
		Db:   Init(),
	}

	return s
}
