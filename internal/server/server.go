package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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

func (s *Server) Start(addr string) error {
	if err := s.Echo.Start(":" + addr); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("start echo: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown echo: %w", err)
	}

	return nil
}
