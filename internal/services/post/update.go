package post

import (
	"echo-demo-project/internal/models"
	"echo-demo-project/internal/requests"
)

func (postService *Service) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.DB.Save(post)
}
