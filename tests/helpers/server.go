package helpers

import (
	"log"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *server.Server {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Failed to parse envs:", err.Error())
	}

	s := &server.Server{
		Echo:   echo.New(),
		DB:     db,
		Config: &cfg,
	}

	return s
}
