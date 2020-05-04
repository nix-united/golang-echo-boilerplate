package seeders

import (
	"echo-demo-project/server/models"

	"github.com/jinzhu/gorm"
)

type UserSeeder struct {
	DB *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{DB: db}
}

func (userSeeder *UserSeeder) SetUsers() {
	users := map[int]map[string]string{
		1: {
			"name":     "user1",
			"password": "password1",
		},
		2: {
			"name":     "user2",
			"password": "password2",
		},
	}

	for key, value := range users {
		user := models.User{}
		userSeeder.DB.First(&user, key)

		if user.ID == 0 {
			user.ID = uint(key)
			user.Name = value["name"]
			user.Password = value["password"]
			userSeeder.DB.Create(&user)
		}
	}
}
