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

func (s Service) GetPosts() ([]models.Post, error) {
	posts, err := s.postRepository.GetPosts()
	if err != nil {
		return nil, fmt.Errorf("get posts from repository: %w", err)
	}

	return posts, nil
}

func (s Service) GetPost(id int) (models.Post, error) {
	post, err := s.postRepository.GetPost(id)
	if err != nil {
		return models.Post{}, fmt.Errorf("get post from repository: %w", err)
	}

	return post, nil
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
