package repositories

import (
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

func (r PostRepository) GetPost(post *models.Post, id int) {
	r.db.Where("id = ? ", id).Find(post)
}

func (r PostRepository) Update(post *models.Post) {
	r.db.Save(post)
}

func (r PostRepository) Delete(post *models.Post) {
	r.db.Delete(post)
}
