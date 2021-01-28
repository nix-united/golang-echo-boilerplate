package server

import (
	"echo-demo-project/server/db"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
	DB   *gorm.DB
}

func NewServer() *Server {
	return &Server{
		Echo: echo.New(),
		DB:   db.Init(),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
