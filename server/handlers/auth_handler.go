package handlers

import (
	"echo-demo-project/server"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"github.com/labstack/echo"
	"net/http"
)

type AuthHandler struct {
	server *server.Server
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}
	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.Db)
	userRepository.GetUser(&user, loginRequest)

	if user.ID == 0 {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}
	tokenService := services.NewTokenService()
	token, err := tokenService.CreateToken(&user)

	if err != nil {
		return err
	}
	return responses.SuccessResponse(c, map[string]string{
		"token": token,
	})
}
