package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	"github.com/labstack/echo/v4"
)

//go:generate go tool mockgen -source=$GOFILE -destination=auth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type authService interface {
	GenerateToken(ctx context.Context, request *requests.LoginRequest) (*responses.LoginResponse, error)
	RefreshToken(ctx context.Context, request *requests.RefreshRequest) (*responses.LoginResponse, error)
}

type AuthHandler struct {
	authService authService
}

func NewAuthHandler(authService authService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
//	@Failure		401		{object}	responses.ErrorResponse
//	@Router			/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var request requests.LoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request", http.StatusBadRequest))
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Required fields are empty or not valid", http.StatusBadRequest))
	}

	response, err := h.authService.GenerateToken(c.Request().Context(), &request)
	switch {
	case errors.Is(err, models.ErrUserNotFound), errors.Is(err, models.ErrInvalidPassword):
		return c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid credentials", http.StatusUnauthorized))
	case err != nil:
		return c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Internal Server Error", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response)
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
//	@Failure		401		{object}	responses.ErrorResponse
//	@Router			/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var request requests.RefreshRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request", http.StatusBadRequest))
	}

	response, err := h.authService.RefreshToken(c.Request().Context(), &request)
	switch {
	case errors.Is(err, models.ErrUserNotFound), errors.Is(err, models.ErrInvalidAuthToken):
		return c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Unauthorized", http.StatusUnauthorized))
	case err != nil:
		return c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Internal Server Error", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response)
}
