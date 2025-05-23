package handlers

import (
	"errors"
	"net/http"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	s "github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"

	"github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	server *s.Server
}

func NewRegisterHandler(server *s.Server) *RegisterHandler {
	return &RegisterHandler{server: server}
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
func (registerHandler *RegisterHandler) Register(c echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}

	if err := registerRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	userRepository := repositories.NewUserRepository(registerHandler.server.DB)

	_, err := userRepository.GetUserByEmail(c.Request().Context(), registerRequest.Email)
	if err == nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "User already exists")
	} else if !errors.Is(err, models.ErrUserNotFound) {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Server error")
	}

	userService := user.NewUserService(registerHandler.server.DB)
	if err := userService.Register(registerRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Server error")
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully created")
}
