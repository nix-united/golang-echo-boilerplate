package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc"
	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/builders"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

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
	idTokenVerifier *oidc.IDTokenVerifier
	tokenService    token.ServiceWrapper
	userRepository  userRepository
}

func NewService(idTokenVerifier *oidc.IDTokenVerifier, tokenService token.ServiceWrapper, userRepository userRepository) *Service {
	return &Service{userRepository: userRepository, tokenService: tokenService, idTokenVerifier: idTokenVerifier}
}

func (s *Service) Register(ctx context.Context, request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
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

func (s *Service) GoogleOAuth(ctx context.Context, token string) (string, string, int64, error) {
	payload, err := s.idTokenVerifier.Verify(ctx, token)
	if err != nil {
		return "", "", 0, fmt.Errorf("verify google token: %w", err)
	}

	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	err = payload.Claims(&claims)
	if err != nil {
		return "", "", 0, fmt.Errorf("extract claims: %w", err)
	}

	if claims.Email == "" {
		return "", "", 0, fmt.Errorf("email is empty")
	}

	user, err := s.userRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotFound) {
			return "", "", 0, fmt.Errorf("get user: %w", err)
		}

		user = *builders.NewUserBuilder().
			SetEmail(claims.Email).
			SetName(claims.Name).
			Build()

		oAuthProvider := models.OAuthProviders{
			UserID:   user.ID,
			Provider: models.GOOGLE,
			Token:    token,
		}

		err = s.userRepository.CreateUserAndOAuthProvider(ctx, &user, &oAuthProvider)
		if err != nil {
			return "", "", 0, fmt.Errorf("create user and oauth provider: %w", err)
		}
	}

	accessToken, i, err := s.tokenService.CreateAccessToken(&user)
	if err != nil {
		return "", "", 0, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := s.tokenService.CreateRefreshToken(&user)
	if err != nil {
		return "", "", 0, fmt.Errorf("create refresh token: %w", err)
	}

	return accessToken, refreshToken, i, nil
}
