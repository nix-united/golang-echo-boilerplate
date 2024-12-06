package token

import (
	"echo-demo-project/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (tokenService *Service) CreateAccessToken(user *models.User) (t string, expired int64, err error) {
	exp := time.Now().Add(time.Hour * ExpireCount)
	claims := &JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
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
