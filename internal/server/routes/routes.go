package routes

import (
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	s "github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/middleware"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/post"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(tracer slogx.TraceStarter, server *s.Server) {
	userRepository := repositories.NewUserRepository(server.DB)
	userService := user.NewService(userRepository)

	postRepository := repositories.NewPostRepository(server.DB)
	postService := post.NewService(postRepository)

	tokenService := token.NewTokenService(server.Config)

	postHandler := handlers.NewPostHandlers(postService)
	authHandler := handlers.NewAuthHandler(server.Config.Auth.RefreshSecret, userService, tokenService)
	registerHandler := handlers.NewRegisterHandler(userService)

	server.Echo.Use(middleware.NewRequestLogger(tracer))

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	r := server.Echo.Group("", middleware.NewRequestDebugger())

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
