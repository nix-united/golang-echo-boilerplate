package post

import (
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
)

type postRepository interface {
	Create(post *models.Post) error
	GetPosts() ([]models.Post, error)
	GetPost(id int) (models.Post, error)
	Update(post *models.Post) error
	Delete(post *models.Post) error
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

func (s Service) GetPosts(posts *[]models.Post) {
	s.postRepository.GetPosts()
}

func (s Service) GetPost(post *models.Post, id int) {
	s.postRepository.GetPost(id)
}

func (s Service) Update(post *models.Post, updatePostRequest requests.UpdatePostRequest) error {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title

	if err := s.postRepository.Update(post); err != nil {
		return fmt.Errorf("update post in repository: %w", err)
	}

	return nil
}

func (s Service) Delete(post *models.Post) error {
	if err := s.postRepository.Delete(post); err != nil {
		return fmt.Errorf("delete post in repository: %w", err)
	}

	return nil
}
