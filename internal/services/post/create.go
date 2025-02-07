package post

import "github.com/nix-united/golang-echo-boilerplate/internal/models"

func (postService *Service) Create(post *models.Post) {
	postService.DB.Create(post)
}
