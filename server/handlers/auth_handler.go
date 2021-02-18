package handlers

import (
	"echo-demo-project/models"
	"echo-demo-project/repositories"
	"echo-demo-project/requests"
	"echo-demo-project/responses"
	s "echo-demo-project/server"
	tokenservice "echo-demo-project/services/token"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	server         *s.Server
	userRepository *repositories.UserRepository
	tokenService   *tokenservice.Service
}

func NewAuthHandler(server *s.Server) *AuthHandler {
	return &AuthHandler{
		server:         server,
		userRepository: repositories.NewUserRepository(server.DB),
		tokenService:   tokenservice.NewTokenService(server),
	}
}

// Login godoc
// @Summary Authenticate a user
// @Description Perform user login
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.LoginRequest true "User's credentials"
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
// @Router /login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}

	if err := loginRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	user := models.User{}
	authHandler.userRepository.GetUserByEmail(&user, loginRequest.Email)
	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, refreshToken, exp, err := authHandler.tokenService.GenerateTokenPair(&user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Perform refresh access token
// @ID user-refresh
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RefreshRequest true "Refresh token"
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
// @Router /refresh [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	claims, err := authHandler.tokenService.ParseToken(refreshRequest.Token,
		authHandler.server.Config.Auth.RefreshSecret)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Not authorized")
	}

	if authHandler.tokenService.ValidateToken(claims, true) != nil {
		return responses.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
	}

	user := new(models.User)
	authHandler.userRepository.GetUser(user, int(claims.ID))
	if user.ID == 0 {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	accessToken, refreshToken, exp, err := authHandler.tokenService.GenerateTokenPair(user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}
