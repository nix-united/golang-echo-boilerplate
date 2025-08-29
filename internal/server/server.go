package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
}

func NewServer(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}

func (s *Server) Start(addr string) error {
	if err := s.echo.Start(":" + addr); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("start echo: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown echo: %w", err)
	}

	return nil
}
