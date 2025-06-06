package repositories

import (
	"errors"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return PostRepository{db: db}
}

func (r PostRepository) Create(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return fmt.Errorf("execute insert post query: %w", err)
	}

	return nil
}

func (r PostRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Preload("User").Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("execute select posts query: %w", err)
	}

	return posts, nil
}

func (r PostRepository) GetPost(id int) (models.Post, error) {
	var post models.Post
	err := r.db.Where("id = ?", id).Take(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Post{}, fmt.Errorf("post not found: %w", err)
	} else if err != nil {
		return models.Post{}, fmt.Errorf("execute select post by id query: %w", err)
	}

	return post, nil
}

func (r PostRepository) Update(post *models.Post) error {
	if err := r.db.Save(post).Error; err != nil {
		return fmt.Errorf("execute update post query: %w", err)
	}

	return nil
}

func (r PostRepository) Delete(post *models.Post) error {
	if err := r.db.Delete(post).Error; err != nil {
		return fmt.Errorf("execute delete post query: %w", err)
	}

	return nil
}
