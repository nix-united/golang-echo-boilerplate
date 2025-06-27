package routes

import (
	"net/http"
	"time"

	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	s "github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/middleware"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/auth"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/post"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoswagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(tracer *slogx.TraceStarter, server *s.Server) {
	userRepository := repositories.NewUserRepository(server.DB)
	userService := user.NewService(userRepository)

	postRepository := repositories.NewPostRepository(server.DB)
	postService := post.NewService(postRepository)

	tokenService := token.NewService(
		time.Now,
		server.Config.Auth.AccessTokenDuration,
		server.Config.Auth.RefreshTokenDuration,
		[]byte(server.Config.Auth.AccessSecret),
		[]byte(server.Config.Auth.RefreshSecret),
	)

	authService := auth.NewService(userService, tokenService)

	postHandler := handlers.NewPostHandlers(postService)
	authHandler := handlers.NewAuthHandler(authService)
	registerHandler := handlers.NewRegisterHandler(userService)

	server.Echo.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	server.Echo.GET("/swagger/*", echoswagger.WrapHandler)

	secureAPI := server.Echo.Group("", middleware.NewRequestLogger(tracer))

	secureAPI.POST("/login", authHandler.Login)
	secureAPI.POST("/register", registerHandler.Register)
	secureAPI.POST("/refresh", authHandler.RefreshToken)

	authorizedAPI := server.Echo.Group(
		"",
		middleware.NewRequestLogger(tracer),
		middleware.NewRequestDebugger(),
		echojwt.WithConfig(echojwt.Config{
			NewClaimsFunc: func(_ echo.Context) jwt.Claims {
				return new(token.JwtCustomClaims)
			},
			SigningKey: []byte(server.Config.Auth.AccessSecret),
		}),
	)

	authorizedAPI.GET("/posts", postHandler.GetPosts)
	authorizedAPI.POST("/posts", postHandler.CreatePost)
	authorizedAPI.DELETE("/posts/:id", postHandler.DeletePost)
	authorizedAPI.PUT("/posts/:id", postHandler.UpdatePost)
}
