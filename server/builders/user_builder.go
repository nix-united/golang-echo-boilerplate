package builders

import "echo-demo-project/models"

type UserBuilder struct {
	email    string
	name     string
	password string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (userBuilder *UserBuilder) SetEmail(email string) (u *UserBuilder) {
	userBuilder.email = email
	return userBuilder
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
		Email:    userBuilder.email,
		Name:     userBuilder.name,
		Password: userBuilder.password,
	}

	return user
}
