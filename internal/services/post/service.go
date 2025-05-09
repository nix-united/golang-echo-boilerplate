package post

import (
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
)

type postRepository interface {
	Create(post *models.Post) error
	GetPosts() ([]models.Post, error)
	GetPost(post *models.Post, id int)
	Update(post *models.Post)
	Delete(post *models.Post)
}

type Service struct {
	postRepository postRepository
}

func NewPostService(postRepository postRepository) Service {
	return Service{postRepository: postRepository}
}

func (s Service) Create(post *models.Post) error {
	if err := s.postRepository.Create(post); err != nil {
		return fmt.Errorf("create post in repository: %w", err)
	}

	return nil
}

func (s Service) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title

	s.postRepository.Update(post)
}

func (s Service) Delete(post *models.Post) {
	s.postRepository.Delete(post)
}
