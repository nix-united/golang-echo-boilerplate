package token

import (
	"echo-demo-project/config"
	"echo-demo-project/models"
	"fmt"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60

type JwtCustomClaims struct {
	ID uint `json:"id"`
	jwtGo.StandardClaims
}

type ServiceWrapper interface {
	GenerateTokenPair(user *models.User) (accessToken, refreshToken string, exp int64, err error)
	ParseToken(tokenString, secret string) (claims *jwtGo.MapClaims, err error)
}

type Service struct {
	config *config.Config
}

func NewTokenService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (tokenService *Service) GenerateTokenPair(user *models.User) (
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	if accessToken, exp, err = tokenService.createToken(user.ID, ExpireAccessMinutes,
		tokenService.config.Auth.AccessSecret); err != nil {
		return
	}

	if refreshToken, _, err = tokenService.createToken(user.ID, ExpireRefreshMinutes,
		tokenService.config.Auth.RefreshSecret); err != nil {
		return
	}

	return
}

func (tokenService *Service) ParseToken(tokenString, secret string) (claims *jwtGo.MapClaims, err error) {
	token, err := jwtGo.Parse(tokenString, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, err
}

func (tokenService *Service) createToken(userID uint, expireMinutes int, secret string) (
	token string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
	claims := &JwtCustomClaims{
		ID: userID,
		StandardClaims: jwtGo.StandardClaims{
			ExpiresAt: exp,
		},
	}
	jwtToken := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))

	return
}
