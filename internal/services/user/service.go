package user

import (
	"context"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/builders"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Register(ctx context.Context, request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := builders.NewUserBuilder().
		SetEmail(request.Email).
		SetName(request.Name).
		SetPassword(string(encryptedPassword)).
		Build()

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user in repository: %w", err)
	}

	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by email from repository: %w", err)
	}

	return user, nil
}
