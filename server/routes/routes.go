package routes

import (
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/services"
	"github.com/labstack/echo/middleware"
)

func ConfigureRoutes(server *server.Server) {
	postHandler := handlers.NewPostHandler(server)
	authHandler := handlers.NewAuthHandler(server)

	server.Echo.POST("/login", authHandler.Login())

	r := server.Echo.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:     &services.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))

	r.GET("/posts", postHandler.GetPosts())
	r.POST("/posts", postHandler.CreatePost())
	r.DELETE("/posts/:id", postHandler.DeletePost())
	r.PUT("/posts/:id", postHandler.UpdatePost())
}
