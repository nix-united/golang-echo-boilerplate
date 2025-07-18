package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go tool mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userService interface {
	GetByID(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type tokenService interface {
	CreateAccessToken(ctx context.Context, user *models.User) (string, int64, error)
	CreateRefreshToken(ctx context.Context, user *models.User) (string, error)
}

type Service struct {
	refreshSecret []byte
	userService   userService
	tokenService  tokenService
}

func NewService(
	refreshSecret []byte,
	userService userService,
	tokenService tokenService,
) *Service {
	return &Service{
		refreshSecret: refreshSecret,
		userService:   userService,
		tokenService:  tokenService,
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
	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.refreshSecret, nil
	})
	if err != nil {
		return nil, errors.Join(fmt.Errorf("parse token: %w", err), models.ErrInvalidAuthToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.Join(errors.New("missing claims"), models.ErrInvalidAuthToken)
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.Join(errors.New("missing id claim"), models.ErrInvalidAuthToken)
	}

	user, err := s.userService.GetByID(ctx, uint(userID))
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
