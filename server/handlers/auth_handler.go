package handlers

import (
	s "echo-demo-project/server"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	server *s.Server
}

type SuccessLoginData struct {
	Token string `json:"token"`
}

func NewAuthHandler(server *s.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

// Login godoc
// @Summary Authenticate a user
// @Description Perform user login
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.LoginRequest true "User's credentials"
// @Success 200 {object} SuccessLoginData
// @Failure 401 {object} responses.Error
// @Router /login [post]
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
	return responses.SuccessResponse(c, SuccessLoginData{
		Token: token,
	})
}
