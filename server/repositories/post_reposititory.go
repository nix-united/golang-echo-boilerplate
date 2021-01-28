package repositories

import (
	"echo-demo-project/server/models"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (postRepository *PostRepository) GetPosts(posts *[]models.Post) {
	postRepository.DB.Find(posts)
}

func (postRepository *PostRepository) GetPost(post *models.Post, id int) {
	postRepository.DB.Where("id = ? ", id).Find(post)
}
