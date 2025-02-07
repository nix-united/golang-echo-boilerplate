package post

import "github.com/nix-united/golang-echo-boilerplate/internal/models"

func (postService *Service) Delete(post *models.Post) {
	postService.DB.Delete(post)
}
