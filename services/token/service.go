package token

import (
	"echo-demo-project/models"
	"echo-demo-project/repositories"
	s "echo-demo-project/server"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60
const AutoLogoffMinutes = 10

type JwtCustomClaims struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	jwtGo.StandardClaims
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

type ServiceWrapper interface {
	GenerateTokenPair(user *models.User) (accessToken, refreshToken string, exp int64, err error)
	ParseToken(tokenString, secret string) (claims *jwtGo.MapClaims, err error)
	ValidateToken(claims *JwtCustomClaims, isRefresh bool) error
}

type Service struct {
	server *s.Server
}

func NewTokenService(server *s.Server) *Service {
	return &Service{
		server: server,
	}
}

func (tokenService *Service) GenerateTokenPair(user *models.User) (accessToken,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID string
	if accessToken, accessUID, exp, err = tokenService.createToken(user.ID, ExpireAccessMinutes,
		tokenService.server.Config.Auth.AccessSecret); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = tokenService.createToken(user.ID, ExpireRefreshMinutes,
		tokenService.server.Config.Auth.RefreshSecret); err != nil {
		return
	}

	cacheJSON, err := json.Marshal(CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})
	tokenService.server.Redis.Set(fmt.Sprintf("token-%d", user.ID), string(cacheJSON),
		time.Minute*AutoLogoffMinutes)

	return
}

func (tokenService *Service) ParseToken(tokenString, secret string) (
	claims *JwtCustomClaims,
	err error,
) {
	token, err := jwtGo.ParseWithClaims(tokenString, &JwtCustomClaims{},
		func(token *jwtGo.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (tokenService *Service) ValidateToken(claims *JwtCustomClaims, isRefresh bool) (
	user *models.User,
	err error,
) {
	var g errgroup.Group
	g.Go(func() error {
		cacheJSON, _ := tokenService.server.Redis.Get(fmt.Sprintf("token-%d", claims.ID)).Result()
		cachedTokens := new(CachedTokens)
		err = json.Unmarshal([]byte(cacheJSON), cachedTokens)

		var tokenUID string
		if isRefresh {
			tokenUID = cachedTokens.RefreshUID
		} else {
			tokenUID = cachedTokens.AccessUID
		}

		if err != nil || tokenUID != claims.UID {
			return errors.New("token not found")
		}

		return nil
	})

	g.Go(func() error {
		user = new(models.User)
		userRepository := repositories.NewUserRepository(tokenService.server.DB)
		userRepository.GetUser(user, int(claims.ID))
		if user.ID == 0 {
			return errors.New("user not found")
		}

		return nil
	})

	err = g.Wait()

	return user, err
}

func (tokenService *Service) createToken(userID uint, expireMinutes int, secret string) (token,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
	uid = uuid.New().String()
	claims := &JwtCustomClaims{
		ID:  userID,
		UID: uid,
		StandardClaims: jwtGo.StandardClaims{
			ExpiresAt: exp,
		},
	}
	jwtToken := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))

	return
}
