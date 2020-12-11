package repositories

import (
	"echo-demo-project/models"

	"github.com/jinzhu/gorm"
)

type UserRepositoryQ interface {
	GetUserByEmail(user *models.User, email string)
}

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryQ {
	return &UserRepository{Db: db}
}

func (userRepository *UserRepository) GetUserByEmail(user *models.User, email string) {
	userRepository.Db.Where("email = ?", email).Find(user)
}
