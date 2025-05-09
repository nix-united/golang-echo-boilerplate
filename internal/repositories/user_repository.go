package repositories

import (
	"errors"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, fmt.Errorf("user not found: %w", err)
	} else if err != nil {
		return models.User{}, fmt.Errorf("execute select user by email query: %w", err)
	}

	return user, nil
}
