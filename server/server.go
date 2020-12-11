package server

import (
	"echo-demo-project/config"
	"echo-demo-project/db"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo   *echo.Echo
	Db     *gorm.DB
	Config config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		Db:     db.Init(cfg),
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
