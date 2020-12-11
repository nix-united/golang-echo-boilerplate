package main

import (
	application "echo-demo-project"
	"echo-demo-project/config"
	"echo-demo-project/server"
	"echo-demo-project/server/routes"
	"fmt"
	"log"
	"os"

	"echo-demo-project/docs"
)

// @title Echo Demo App
// @version 1.0
// @description This is a demo version of Echo app.

// @contact.name NIX Solutions
// @contact.url https://www.nixsolutions.com/
// @contact.email ask@nixsolutions.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /
func main() {
	cfg := config.NewConfig()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.ExposePort)

	app := server.NewServer(cfg)

	routes.ConfigureRoutes(app)
	err := app.Start(os.Getenv("PORT"))

	if err != nil {
		log.Fatal("Port already used")
	}

	application.Start(cfg)
}
