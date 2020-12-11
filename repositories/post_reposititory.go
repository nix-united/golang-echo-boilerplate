package repositories

import (
	"echo-demo-project/models"

	"github.com/jinzhu/gorm"
)

type PostRepositoryQ interface {
	GetPosts(posts *[]models.Post)
	GetPost(post *models.Post, id int)
}

type PostRepository struct {
	Db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepositoryQ {
	return &PostRepository{Db: db}
}

func (postRepository *PostRepository) GetPosts(posts *[]models.Post) {
	postRepository.Db.Find(posts)
}

func (postRepository *PostRepository) GetPost(post *models.Post, id int) {
	postRepository.Db.Where("id = ? ", id).Find(post)
}
