package routes

import (
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	s "github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/post"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server) {
	postRepository := repositories.NewPostRepository(server.DB)
	postService := post.NewPostService(postRepository)

	postHandler := handlers.NewPostHandlers(postService)
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	server.Echo.Use(middleware.Logger())

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	fmt.Println(server.Config.Auth.AccessSecret)

	r := server.Echo.Group("")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(server.Config.Auth.AccessSecret),
	}
	r.Use(echojwt.WithConfig(config))

	r.GET("/posts", postHandler.GetPosts)
	r.POST("/posts", postHandler.CreatePost)
	r.DELETE("/posts/:id", postHandler.DeletePost)
	r.PUT("/posts/:id", postHandler.UpdatePost)
}
