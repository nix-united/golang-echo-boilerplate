package server

import (
	"echo-demo-project/server/db"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Server struct {
	Echo *echo.Echo
	Db   *gorm.DB
}

func NewServer() *Server {
	return &Server{
		Echo: echo.New(),
		Db:   db.Init(),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
