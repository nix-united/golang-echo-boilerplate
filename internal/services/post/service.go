package post

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/domain"
	"github.com/nix-united/golang-echo-boilerplate/internal/models"
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

// UpdateByUser checks if user has rights to update a provided post and updates it.
func (s *Service) UpdateByUser(ctx context.Context, request domain.UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepository.GetPost(ctx, request.PostID)
	if err != nil {
		return nil, fmt.Errorf("get stored post from repository: %w", err)
	}

	if post.UserID != request.UserID {
		return nil, models.ErrForbidden
	}

	post.Title = request.Title
	post.Content = request.Content

	if err := s.postRepository.Update(ctx, &post); err != nil {
		return nil, fmt.Errorf("update post in repository: %w", err)
	}

	return &post, nil
}

// DeleteByUser checks if user has rights to delete a post and deletes it.
func (s *Service) DeleteByUser(ctx context.Context, request domain.DeletePostRequest) error {
	post, err := s.postRepository.GetPost(ctx, request.PostID)
	if err != nil {
		return fmt.Errorf("get stored post from repository: %w", err)
	}

	if post.UserID != request.UserID {
		return models.ErrForbidden
	}

	if err := s.postRepository.Delete(ctx, &post); err != nil {
		return fmt.Errorf("delete post in repository: %w", err)
	}

	return nil
}
