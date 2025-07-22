package token

import (
	"context"
	"fmt"
	"time"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

type Service struct {
	now                  func() time.Time
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	accessTokenSecret    []byte
	refreshSecret        []byte
}

func NewService(
	now func() time.Time,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
	accessSecret []byte,
	refreshSecret []byte,
) *Service {
	return &Service{
		now:                  now,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		accessTokenSecret:    accessSecret,
		refreshSecret:        refreshSecret,
	}
}

func (s *Service) CreateAccessToken(_ context.Context, user *models.User) (accessToken string, expires int64, err error) {
	expiresAt := s.now().Add(s.accessTokenDuration)

	claims := &JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString(s.accessTokenSecret)
	if err != nil {
		return "", 0, fmt.Errorf("sign access token: %w", err)
	}

	return accessToken, expiresAt.Unix(), nil
}

func (s *Service) CreateRefreshToken(_ context.Context, user *models.User) (string, error) {
	expiresAt := s.now().Add(s.refreshTokenDuration)

	claims := &JwtCustomRefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.refreshSecret)
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	return signed, nil
}

func (s *Service) ParseAccessToken(_ context.Context, token string) (*JwtCustomClaims, error) {
	claims := new(JwtCustomClaims)
	if err := s.parseToken(token, s.accessTokenSecret, claims); err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}

func (s *Service) ParseRefreshToken(_ context.Context, token string) (*JwtCustomRefreshClaims, error) {
	claims := new(JwtCustomRefreshClaims)
	if err := s.parseToken(token, s.refreshSecret, claims); err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}

func (s *Service) parseToken(token string, secret []byte, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return fmt.Errorf("parse token with claims: %w", err)
	}

	return nil
}
