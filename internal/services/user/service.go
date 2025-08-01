package user

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go tool mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oauthProvider *models.OAuthProviders) error
}

type Service struct {
	userRepository userRepository
}

func NewService(userRepository userRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(ctx context.Context, request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	user := &models.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: string(encryptedPassword),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user in repository: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (models.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id from repository: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by email from repository: %w", err)
	}

	return user, nil
}

func (s *Service) CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oauthProvider *models.OAuthProviders) error {
	err := s.userRepository.CreateUserAndOAuthProvider(ctx, user, oauthProvider)
	if err != nil {
		return fmt.Errorf("create user and oauth provider from repository: %w", err)
	}

	return nil
}
