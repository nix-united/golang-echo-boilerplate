package services

import (
	"echo-demo-project/server/models"
	"github.com/jinzhu/gorm"
)

type PostService struct {
	Db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{Db : db}
}

func (postService *PostService) Create(post *models.Post){
	postService.Db.Create(post)
}

