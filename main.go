package main

import (
	"echo-demo-project/server"
	"echo-demo-project/server/routes"
	"echo-demo-project/server/validation"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"

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
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("EXPOSE_PORT"))

	app := server.NewServer()
	app.Echo.Validator = validation.NewCustomValidator(validator.New())

	routes.ConfigureRoutes(app)
	err = app.Start(os.Getenv("PORT"))

	if err != nil {
		log.Fatal("Port already used")
	}
}
