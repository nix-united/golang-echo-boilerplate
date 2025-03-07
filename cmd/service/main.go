package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nix-united/golang-echo-boilerplate/docs"
	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/db"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/routes"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
		os.Exit(1)
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

	dbConnection, err := db.NewConnection(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db connection: %w", err)
	}

	app := server.NewServer(echo.New(), dbConnection, &cfg)

	routes.ConfigureRoutes(app)

	err = app.Start(cfg.HTTP.Port)
	if err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}
