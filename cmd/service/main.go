package main

import (
	"fmt"
	"log/slog"

	application "github.com/nix-united/golang-echo-boilerplate"
	"github.com/nix-united/golang-echo-boilerplate/docs"
	"github.com/nix-united/golang-echo-boilerplate/internal/config"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

//	@title			Echo Demo App
//	@version		1.0
//	@description	This is a demo version of Echo app.

//	@contact.name	NIX Solutions
//	@contact.url	https://www.nixsolutions.com/
//	@contact.email	ask@nixsolutions.com

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/
func main() {
	if err := run(); err != nil {
		slog.Error("Service run error", "err", err.Error())
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	application.Start(&cfg)

	return nil
}
