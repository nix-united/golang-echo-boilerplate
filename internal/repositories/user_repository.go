package repositories

import (
	"context"
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

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("execute insert user query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Join(models.ErrUserNotFound, err)
	} else if err != nil {
		return models.User{}, fmt.Errorf("execute select user by id query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Join(models.ErrUserNotFound, err)
	} else if err != nil {
		return models.User{}, fmt.Errorf("execute select user by email query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oAuthProvider *models.OAuthProviders) error {
	tx := r.db.Begin()

	committed := false

	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("insert user (tx): %w", err)
	}

	oAuthProvider.UserID = user.ID

	if err := tx.WithContext(ctx).Create(oAuthProvider).Error; err != nil {
		return fmt.Errorf("insert oauthprovider (tx): %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	return nil
}
