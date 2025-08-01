package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"github.com/coreos/go-oidc"
)

type Service struct {
	idTokenVerifier *oidc.IDTokenVerifier
	tokenService    tokenService
	userService     userService
}

type userService interface {
	CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oAuthProvider *models.OAuthProviders) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type tokenService interface {
	CreateAccessToken(ctx context.Context, user *models.User) (string, int64, error)
	CreateRefreshToken(ctx context.Context, user *models.User) (string, error)
}

func NewService(idTokenVerifier *oidc.IDTokenVerifier, tokenService tokenService, userService userService) *Service {
	return &Service{idTokenVerifier: idTokenVerifier, tokenService: tokenService, userService: userService}
}

func (s Service) GoogleOAuth(ctx context.Context, token string) (accessToken, refreshToken string, exp int64, err error) {
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

	user, err := s.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotFound) {
			return "", "", 0, fmt.Errorf("get user: %w", err)
		}

		user = models.User{
			Email: claims.Email,
			Name:  claims.Name,
		}

		oAuthProvider := models.OAuthProviders{
			UserID:   user.ID,
			Provider: models.GOOGLE,
			Token:    token,
		}

		err = s.userService.CreateUserAndOAuthProvider(ctx, &user, &oAuthProvider)
		if err != nil {
			return "", "", 0, fmt.Errorf("create user and oauth provider: %w", err)
		}
	}

	accessToken, exp, err = s.tokenService.CreateAccessToken(ctx, &user)
	if err != nil {
		return "", "", 0, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err = s.tokenService.CreateRefreshToken(ctx, &user)
	if err != nil {
		return "", "", 0, fmt.Errorf("create refresh token: %w", err)
	}

	return accessToken, refreshToken, exp, nil
}
