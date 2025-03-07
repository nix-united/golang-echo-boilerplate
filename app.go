package application

import (
	"log"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/db"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/routes"

	"github.com/labstack/echo/v4"
)

func Start(cfg *config.Config) {
	dbConnection, err := db.NewConnection(cfg.DB)
	if err != nil {
		log.Fatal("DB connection error: " + err.Error())
	}

	app := server.NewServer(echo.New(), dbConnection, cfg)

	routes.ConfigureRoutes(app)

	err = app.Start(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used: " + err.Error())
	}
}
