package services

import (
	"echo-demo-project/server/models"
	"echo-demo-project/server/requests"
	"github.com/jinzhu/gorm"
)

type PostService struct {
	Db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{Db: db}
}

func (postService *PostService) Create(post *models.Post) {
	postService.Db.Create(post)
}

func (postService *PostService) Delete(post *models.Post) {
	postService.Db.Delete(post)
}

func (postService *PostService) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.Db.Save(post)
}
