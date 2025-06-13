package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	s "github.com/nix-united/golang-echo-boilerplate/internal/server"
	tokenservice "github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type userGetter interface {
	GetByID(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type AuthHandler struct {
	userGetter userGetter
	server     *s.Server
}

func NewAuthHandler(userGetter userGetter, server *s.Server) *AuthHandler {
	return &AuthHandler{
		server:     server,
		userGetter: userGetter,
	}
}

// Login godoc
//
//	@Summary		Authenticate a user
//	@Description	Perform user login
//	@ID				user-login
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.LoginRequest	true	"User's credentials"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}

	if err := loginRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	user, err := h.userGetter.GetUserByEmail(c.Request().Context(), loginRequest.Email)
	if errors.Is(err, models.ErrUserNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "User with such email not found")
	} else if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	tokenService := tokenservice.NewTokenService(h.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Perform refresh access token
//	@ID				user-refresh
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	token, err := jwtGo.Parse(refreshRequest.Token, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.server.Config.Auth.RefreshSecret), nil
	})

	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	claims, ok := token.Claims.(jwtGo.MapClaims)
	if !ok && !token.Valid {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
	}

	user, err := h.userGetter.GetByID(c.Request().Context(), uint(claims["id"].(float64)))
	if errors.Is(err, models.ErrUserNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "User with such email not found")
	} else if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
	}

	tokenService := tokenservice.NewTokenService(h.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}

	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}
