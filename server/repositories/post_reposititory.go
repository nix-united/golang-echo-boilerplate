package repositories

import (
	"echo-demo-project/server/models"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	Db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{Db: db}
}

func (postRepository *PostRepository) GetPosts(posts *[]models.Post) {
	postRepository.Db.Find(posts)
}

func (postRepository *PostRepository) GetPost(post *models.Post, id int) {
	postRepository.Db.Where("id = ? ", id).Find(post)
}
