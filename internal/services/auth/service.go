package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go tool mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userService interface {
	GetByID(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type tokenService interface {
	ParseRefreshToken(ctx context.Context, token string) (*token.JwtCustomRefreshClaims, error)
	CreateAccessToken(ctx context.Context, user *models.User) (string, int64, error)
	CreateRefreshToken(ctx context.Context, user *models.User) (string, error)
}

type Service struct {
	userService  userService
	tokenService tokenService
}

func NewService(userService userService, tokenService tokenService) *Service {
	return &Service{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *Service) GenerateToken(ctx context.Context, request *requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := s.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.Join(fmt.Errorf("compare hash and passowrd: %w", err), models.ErrInvalidPassword)
	}

	accessToken, exp, err := s.tokenService.CreateAccessToken(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := s.tokenService.CreateRefreshToken(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return response, nil
}

func (s *Service) RefreshToken(ctx context.Context, request *requests.RefreshRequest) (*responses.LoginResponse, error) {
	claims, err := s.tokenService.ParseRefreshToken(ctx, request.Token)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("parse token: %w", err), models.ErrInvalidAuthToken)
	}

	user, err := s.userService.GetByID(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	accessToken, exp, err := s.tokenService.CreateAccessToken(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := s.tokenService.CreateRefreshToken(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return response, nil
}
