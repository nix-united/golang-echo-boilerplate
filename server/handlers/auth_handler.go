package handlers

import (
	s "echo-demo-project/server"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"fmt"
	"net/http"
	"os"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	server *s.Server
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
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
// @Router /login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}
	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.Db)
	userRepository.GetUserByName(&user, loginRequest.Name)

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	tokenService := services.NewTokenService()
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

// Refresh godoc
// @Summary Refresh access token
// @Description Perform refresh access token
// @ID user-refresh
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RefreshRequest true "Access token"
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
// @Router /refresh [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	token, err := jwtGo.Parse(refreshRequest.Token, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	claims, ok := token.Claims.(jwtGo.MapClaims)
	if !ok && !token.Valid {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
	}

	user := new(models.User)
	authHandler.server.Db.First(&user, int(claims["id"].(float64)))

	if user.ID == 0 {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	tokenService := services.NewTokenService()
	accessToken, exp, err := tokenService.CreateAccessToken(user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}
