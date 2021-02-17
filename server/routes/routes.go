package routes

import (
	"echo-demo-project/middleware"
	s "echo-demo-project/server"
	"echo-demo-project/server/handlers"

	echoMW "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server) {
	postHandler := handlers.NewPostHandlers(server)
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	server.Echo.Use(echoMW.Logger())

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	authMW := middleware.JWT(server.Config.Auth.AccessSecret)
	apiProtected := server.Echo.Group("")
	apiProtected.Use(authMW)

	apiProtected.GET("/posts", postHandler.GetPosts)
	apiProtected.POST("/posts", postHandler.CreatePost)
	apiProtected.DELETE("/posts/:id", postHandler.DeletePost)
	apiProtected.PUT("/posts/:id", postHandler.UpdatePost)
}
