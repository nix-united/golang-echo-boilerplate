package main

import (
	"fmt"

	application "github.com/nix-united/golang-echo-boilerplate"
	"github.com/nix-united/golang-echo-boilerplate/docs"
	"github.com/nix-united/golang-echo-boilerplate/internal/config"
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
	cfg := config.NewConfig()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	application.Start(cfg)
}
