package middleware

import (
	"echo-demo-project/responses"
	s "echo-demo-project/server"
	tokenService "echo-demo-project/services/token"
	"fmt"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
)

// Default middleware to check the token.
func JWT(secret string) echo.MiddlewareFunc {
	config := echoMW.JWTConfig{
		Claims:     &tokenService.JwtCustomClaims{},
		SigningKey: []byte(secret),
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Not authorized",
			}
		},
	}

	return echoMW.JWTWithConfig(config)
}

// Middleware for additional steps:
// 1. Check the user exists in DB
// 2. Check the token info exists in Redis
// 3. Add the user DB data to Context
// 4. Prolong the Redis TTL of the current token pair
func ValidateJWT(server *s.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwtGo.Token)
			claims := token.Claims.(*tokenService.JwtCustomClaims)

			user, err := tokenService.NewTokenService(server).ValidateToken(claims, false)
			if err != nil {
				return responses.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
			}

			c.Set("currentUser", user)

			go func() {
				server.Redis.Expire(fmt.Sprintf("token-%d", claims.ID),
					time.Minute*tokenService.AutoLogoffMinutes)
			}()

			return next(c)
		}
	}
}
