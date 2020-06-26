package handlers

import (
	s "echo-demo-project/server"
	"echo-demo-project/server/builders"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	server *s.Server
}

func NewRegisterHandler(server *s.Server) *RegisterHandler {
	return &RegisterHandler{server: server}
}

// Register godoc
// @Summary Register
// @Description New user registration
// @ID user-register
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RegisterRequest true "User's email, user's password"
// @Success 200 {string} string "User successfully created"
// @Failure 400 {object} responses.Error
// @Router /register [post]
func (registerHandler *RegisterHandler) Register(c echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}
	if err := c.Validate(registerRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	existUser := models.User{}
	userRepository := repositories.NewUserRepository(registerHandler.server.Db)
	userRepository.GetUserByName(&existUser, registerRequest.Name)

	if existUser.ID != 0 {
		return responses.ErrorResponse(c, http.StatusBadRequest, "User already exist")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(registerRequest.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Server error")
	}

	user := builders.NewUserBuilder().SetName(registerRequest.Name).
		SetPassword(string(encryptedPassword)).
		Build()

	userService := services.NewUserService(registerHandler.server.Db)
	userService.Create(&user)

	return responses.SuccessResponse(c, "User successfully created")
}
