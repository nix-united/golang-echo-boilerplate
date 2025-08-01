package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nix-united/golang-echo-boilerplate/docs"
	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/db"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/routes"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

const shutdownTimeout = 20 * time.Second

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

	if err := slogx.Init(cfg.Logger); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	gormDB, err := db.NewGormDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db connection: %w", err)
	}

	app := server.NewServer(echo.New(), gormDB, &cfg)

	err = routes.ConfigureRoutes(slogx.NewTraceStarter(uuid.NewV7), app)
	if err != nil {
		return fmt.Errorf("configure routes: %w", err)
	}

	go func() {
		if err = app.Start(cfg.HTTP.Port); err != nil {
			slog.Error("Server error", "err", err.Error())
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownChannel

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	dbConnection, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get db connection: %w", err)
	}

	if err := dbConnection.Close(); err != nil {
		return fmt.Errorf("close db connection: %w", err)
	}

	return nil
}
