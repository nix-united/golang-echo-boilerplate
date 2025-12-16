package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/nix-united/golang-echo-boilerplate/docs"
	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/db"
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/middleware"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/routes"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/auth"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/oauth"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/post"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"

	"github.com/caarlos0/env/v11"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
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
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
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

	userRepository := repositories.NewUserRepository(gormDB)
	userService := user.NewService(userRepository)

	postRepository := repositories.NewPostRepository(gormDB)
	postService := post.NewService(postRepository)

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("oidc.NewProvider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OAuth.ClientID})

	tokenService := token.NewService(
		time.Now,
		cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
		[]byte(cfg.Auth.AccessSecret),
		[]byte(cfg.Auth.RefreshSecret),
	)

	authService := auth.NewService(userService, tokenService)
	oAuthService := oauth.NewService(verifier, tokenService, userService)

	postHandler := handlers.NewPostHandlers(postService)
	authHandler := handlers.NewAuthHandler(authService)
	oAuthHandler := handlers.NewOAuthHandler(oAuthService)
	registerHandler := handlers.NewRegisterHandler(userService)

	// Configure middleware with the custom claims type
	echoJWTConfig := echojwt.Config{
		NewClaimsFunc: func(echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.Auth.AccessSecret),
	}

	echoJWTMiddleware := echojwt.WithConfig(echoJWTConfig)

	engine := routes.ConfigureRoutes(routes.Handlers{
		PostHandler:               postHandler,
		AuthHandler:               authHandler,
		OAuthHandler:              oAuthHandler,
		RegisterHandler:           registerHandler,
		EchoJWTMiddleware:         echoJWTMiddleware,
		RequestLoggerMiddleware:   middleware.NewRequestLogger(slogx.NewTraceStarter(uuid.NewV7)),
		RequestDebuggerMiddleware: middleware.NewRequestDebugger(),
	})
	if err != nil {
		return fmt.Errorf("configure routes: %w", err)
	}

	app := server.NewServer(engine)
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
