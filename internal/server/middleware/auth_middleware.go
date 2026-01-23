package middleware

import (
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// authContextKey is a key to use in context to propagate user data between middlewares.
const authContextKey = "user"

func NewAuthMiddleware(accessSecret string) echo.MiddlewareFunc {
	echoJWTConfig := echojwt.Config{
		NewClaimsFunc: func(echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(accessSecret),
		SuccessHandler: func(c echo.Context) {
			user, ok := c.Get(authContextKey).(*jwt.Token)
			if !ok {
				return
			}

			claims, ok := user.Claims.(*token.JwtCustomClaims)
			if !ok {
				return
			}

			// Enrich logs and context execution with user ID.
			ctx := c.Request().Context()
			ctx = slogx.ContextWithUserID(ctx, claims.ID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
		ContextKey: authContextKey,
	}

	return echojwt.WithConfig(echoJWTConfig)
}
