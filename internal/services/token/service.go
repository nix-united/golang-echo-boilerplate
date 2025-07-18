package token

import (
	"context"
	"time"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const ExpireCount = 2
const ExpireRefreshCount = 168

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

type ServiceWrapper interface {
	CreateAccessToken(user *models.User) (accessToken string, exp int64, err error)
	CreateRefreshToken(user *models.User) (t string, err error)
}

type Service struct {
	config *config.Config
}

func NewTokenService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (tokenService *Service) CreateAccessToken(_ context.Context, user *models.User) (t string, expired int64, err error) {
	exp := time.Now().Add(time.Hour * ExpireCount)
	claims := &JwtCustomClaims{
		user.Name,
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	expired = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(tokenService.config.Auth.AccessSecret))
	if err != nil {
		return
	}

	return
}

func (tokenService *Service) CreateRefreshToken(_ context.Context, user *models.User) (t string, err error) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * ExpireRefreshCount)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	rt, err := refreshToken.SignedString([]byte(tokenService.config.Auth.RefreshSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}
