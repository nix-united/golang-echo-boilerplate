package post

import "echo-demo-project/models"

func (postService *Service) Delete(post *models.Post) {
	postService.Db.Delete(post)
}
