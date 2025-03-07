package server

import (
	"github.com/nix-united/golang-echo-boilerplate/internal/config"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
}

func NewServer(
	Echo *echo.Echo,
	DB *gorm.DB,
	Config *config.Config,
) *Server {
	return &Server{
		Echo:   Echo,
		DB:     DB,
		Config: Config,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
