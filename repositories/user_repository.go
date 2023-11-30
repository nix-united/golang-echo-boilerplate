package repositories

import (
	"echo-demo-project/models"

	"gorm.io/gorm"
)

type UserRepositoryQ interface {
	GetUserByEmail(user *models.User, email string)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepository *UserRepository) GetUserByEmail(user *models.User, email string) {
	userRepository.DB.Where("email = ?", email).Find(user)
}
