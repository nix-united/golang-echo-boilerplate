package builders

import "echo-demo-project/server/models"

type UserBuilder struct {
	name     string
	password string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (userBuilder *UserBuilder) SetName(name string) (u *UserBuilder) {
	userBuilder.name = name
	return userBuilder
}

func (userBuilder *UserBuilder) SetPassword(password string) (u *UserBuilder) {
	userBuilder.password = password
	return userBuilder
}

func (userBuilder *UserBuilder) Build() models.User {
	user := models.User{
		Name:     userBuilder.name,
		Password: userBuilder.password,
	}

	return user
}
