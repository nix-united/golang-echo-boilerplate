package post

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
)

//go:generate go tool mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetPosts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id uint) (models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, post *models.Post) error
}

type Service struct {
	postRepository postRepository
}

func NewService(postRepository postRepository) *Service {
	return &Service{postRepository: postRepository}
}

func (s *Service) Create(ctx context.Context, post *models.Post) error {
	if err := s.postRepository.Create(ctx, post); err != nil {
		return fmt.Errorf("create post in repository: %w", err)
	}

	return nil
}

func (s *Service) GetPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := s.postRepository.GetPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get posts from repository: %w", err)
	}

	return posts, nil
}

func (s *Service) GetPost(ctx context.Context, id uint) (models.Post, error) {
	post, err := s.postRepository.GetPost(ctx, id)
	if err != nil {
		return models.Post{}, fmt.Errorf("get post from repository: %w", err)
	}

	return post, nil
}

func (s *Service) Update(ctx context.Context, post *models.Post, updatePostRequest requests.UpdatePostRequest) error {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title

	if err := s.postRepository.Update(ctx, post); err != nil {
		return fmt.Errorf("update post in repository: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, post *models.Post) error {
	if err := s.postRepository.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete post in repository: %w", err)
	}

	return nil
}
