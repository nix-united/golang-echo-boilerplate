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

//go:generate go tool mockgen -source=$GOFILE -destination=register_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRegisterer interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	Register(ctx context.Context, request *requests.RegisterRequest) error
}

type RegisterHandler struct {
	userRegisterer userRegisterer
}

func NewRegisterHandler(userRegisterer userRegisterer) *RegisterHandler {
	return &RegisterHandler{userRegisterer: userRegisterer}
}

// Register godoc
//
//	@Summary		Register
//	@Description	New user registration
//	@ID				user-register
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.RegisterRequest	true	"User's email, user's password"
//	@Success		201		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Router			/register [post]
func (h *RegisterHandler) Register(c echo.Context) error {
	var registerRequest requests.RegisterRequest
	if err := c.Bind(&registerRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request")
	}

	if err := registerRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or invalid")
	}

	_, err := h.userRegisterer.GetUserByEmail(c.Request().Context(), registerRequest.Email)
	if err == nil {
		return responses.ErrorResponse(c, http.StatusConflict, "User already exists")
	} else if !errors.Is(err, models.ErrUserNotFound) {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to check if user exists")
	}

	if err := h.userRegisterer.Register(c.Request().Context(), &registerRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully created")
}
