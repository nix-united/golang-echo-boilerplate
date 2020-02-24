package repositories

import (
	"echo-demo-project/server/models"
	"echo-demo-project/server/requests"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (userRepository *UserRepository) GetUser(user *models.User, loginRequest *requests.LoginRequest) {
	userRepository.Db.Where("name = ?", loginRequest.Name).
		Where("password = ?", loginRequest.Password).
		Find(user)
}

func (userRepository *UserRepository) GetUserByName(user *models.User, name string) {
	userRepository.Db.Where("name = ?", name).Find(user)
}
