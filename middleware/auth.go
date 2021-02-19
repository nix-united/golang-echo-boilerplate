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

func ValidateJWT(server *s.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwtGo.Token)
			claims := token.Claims.(*tokenService.JwtCustomClaims)

			if tokenService.NewTokenService(server).ValidateToken(claims, false) != nil {
				return responses.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
			}

			server.Redis.Expire(fmt.Sprintf("token-%d", claims.ID),
				time.Minute*tokenService.AutoLogoffMinutes)

			return next(c)
		}
	}
}
