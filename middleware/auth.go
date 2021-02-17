package middleware

import (
	"echo-demo-project/services/token"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
)

func JWT(secret string) echo.MiddlewareFunc {
	config := echoMW.JWTConfig{
		Claims:     &token.JwtCustomClaims{},
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
