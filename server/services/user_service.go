package services

import (
	"echo-demo-project/server/models"
	"github.com/jinzhu/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{Db: db}
}

func (userService *UserService) Create(user *models.User) {
	userService.Db.Create(user)
}
