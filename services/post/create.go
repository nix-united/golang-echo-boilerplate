package post

import "echo-demo-project/models"

func (postService *Service) Create(post *models.Post) {
	postService.Db.Create(post)
}
