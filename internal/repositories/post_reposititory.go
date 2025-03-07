package repositories

import (
	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return PostRepository{db: db}
}

func (r PostRepository) Create(post *models.Post) {
	r.db.Create(post)
}

func (r PostRepository) GetPosts(posts *[]models.Post) {
	r.db.Preload("User").Find(posts)
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
