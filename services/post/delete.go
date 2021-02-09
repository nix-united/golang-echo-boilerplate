package post

import "echo-demo-project/models"

func (postService *Service) Delete(post *models.Post) {
	postService.DB.Delete(post)
}
