package post

import "echo-demo-project/internal/models"

func (postService *Service) Create(post *models.Post) {
	postService.DB.Create(post)
}
