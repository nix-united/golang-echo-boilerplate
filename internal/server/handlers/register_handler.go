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
//	@Success		201		{object}	responses.MessageResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Router			/register [post]
func (h *RegisterHandler) Register(c echo.Context) error {
	var registerRequest requests.RegisterRequest
	if err := c.Bind(&registerRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request", http.StatusBadRequest))
	}

	if err := registerRequest.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Required fields are empty or invalid", http.StatusBadRequest))
	}

	_, err := h.userRegisterer.GetUserByEmail(c.Request().Context(), registerRequest.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, responses.NewErrorResponse("User already exists", http.StatusConflict))
	} else if !errors.Is(err, models.ErrUserNotFound) {
		errorResponse := responses.NewErrorResponse("Failed to check if user exists", http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	if err := h.userRegisterer.Register(c.Request().Context(), &registerRequest); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to register user", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusCreated, responses.NewMessageResponse("User successfully created"))
}
