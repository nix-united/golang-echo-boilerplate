package handlers

import (
	"echo-demo-project/server"
	"echo-demo-project/server/builders"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"github.com/labstack/echo"
	"net/http"
)

type RegisterHandler struct {
	server *server.Server
}

func NewRegisterHandler(server *server.Server) *RegisterHandler {
	return &RegisterHandler{server: server}
}

func (registerHandler *RegisterHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		loginRequest := new(requests.LoginRequest)

		if err := c.Bind(loginRequest); err != nil {
			return err
		}
		if err := c.Validate(loginRequest); err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		}

		existUser := models.User{}
		userRepository := repositories.NewUserRepository(registerHandler.server.Db)
		userRepository.GetUserByName(&existUser, loginRequest.Name)

		if existUser.ID != 0 {
			return responses.ErrorResponse(c, http.StatusBadRequest, "User already exist")
		}

		user := builders.NewUserBuilder().SetName(loginRequest.Name).
			SetPassword(loginRequest.Password).
			Build()

		userService := services.NewUserService(registerHandler.server.Db)
		userService.Create(&user)

		return responses.SuccessResponse(c, "User successfully creat")
	}
}
