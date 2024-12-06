package user

import (
	"echo-demo-project/internal/requests"
	"echo-demo-project/internal/server/builders"

	"golang.org/x/crypto/bcrypt"
)

func (userService *Service) Register(request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := builders.NewUserBuilder().SetEmail(request.Email).
		SetName(request.Name).
		SetPassword(string(encryptedPassword)).
		Build()

	userService.DB.Create(&user)

	return nil
}
