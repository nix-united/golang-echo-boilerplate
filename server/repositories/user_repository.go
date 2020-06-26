package repositories

import (
	"echo-demo-project/server/models"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (userRepository *UserRepository) GetUserByName(user *models.User, name string) {
	userRepository.Db.Where("name = ?", name).Find(user)
}
