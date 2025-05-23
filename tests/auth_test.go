package tests

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/server"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"
	"github.com/nix-united/golang-echo-boilerplate/tests/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestWalkAuth(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		panic(err.Error())
	}

	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/login",
	}
	handlerFunc := func(s *server.Server, c echo.Context) error {
		userRepository := repositories.NewUserRepository(s.DB)
		userService := user.NewUserService(userRepository)
		return handlers.NewAuthHandler(userService, s).Login(c)
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	commonMock := &helpers.QueryMock{
		Query:    "SELECT * FROM `users`  WHERE email = ? AND `users`.`deleted_at` IS NULL",
		QueryArg: []driver.Value{"name@test.com"},
		Reply: helpers.MockReply{
			Columns: []string{"id", "email", "name", "password"},
			Rows: [][]driver.Value{
				{helpers.UserId, "name@test.com", "User Name", encryptedPassword},
			},
		},
	}

	cases := []helpers.TestCase{
		{
			"Auth success",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "password",
				},
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock, commonMock},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "",
			},
		},
		{
			"Login attempt with incorrect password",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "incorrectPassword",
				},
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock, commonMock},
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
		{
			"Login attempt as non-existent user",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "user.not.exists@test.com",
					Password: "password",
				},
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock, commonMock},
			helpers.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "User with such email not found",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			helpers.PrepareDatabaseQueryMocks(test, sqlMock)
			db := helpers.InitGorm(dbMock)
			s := helpers.NewServer(db)

			c, recorder := helpers.PrepareContextFromTestCase(s, test)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				if assert.Equal(t, test.Expected.StatusCode, recorder.Code) {
					if recorder.Code == http.StatusOK {
						assertTokenResponse(t, recorder)
					}
				}
			}
		})
	}
}

func TestWalkRefresh(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		panic(err.Error())
	}

	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/refresh",
	}
	handlerFunc := func(s *server.Server, c echo.Context) error {
		userRepository := repositories.NewUserRepository(s.DB)
		userService := user.NewUserService(userRepository)
		return handlers.NewAuthHandler(userService, s).RefreshToken(c)
	}

	var cfg config.Config
	err = env.Parse(&cfg)
	require.NoError(t, err)

	tokenService := token.NewTokenService(&cfg)

	validUser := models.User{Email: "name@test.com"}
	validUser.ID = helpers.UserId
	validToken, _ := tokenService.CreateRefreshToken(&validUser)

	notExistUser := models.User{Email: "user.not.exists@test.com"}
	notExistUser.ID = helpers.UserId + 1
	notExistToken, _ := tokenService.CreateRefreshToken(&notExistUser)

	invalidToken := validToken[1 : len(validToken)-1]

	cases := []helpers.TestCase{
		{
			"Refresh success",
			request,
			requests.RefreshRequest{
				Token: validToken,
			},
			handlerFunc,
			[]*helpers.QueryMock{
				&helpers.SelectVersionMock,
				{
					Query:    "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1",
					QueryArg: []driver.Value{int64(1)},
					Reply: helpers.MockReply{
						Columns: []string{"id", "name"},
						Rows: [][]driver.Value{
							{helpers.UserId, "User Name"},
						},
					},
				}},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "",
			},
		},
		{
			"Refresh token of non-existent user",
			request,
			requests.RefreshRequest{
				Token: notExistToken,
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock,
				{
					Query:    "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1",
					QueryArg: []driver.Value{int64(2)},
					Reply: helpers.MockReply{
						Columns: []string{"id", "name"},
					},
				},
			},
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "User not found",
			},
		},
		{
			"Refresh invalid token",
			request,
			requests.RefreshRequest{
				Token: invalidToken,
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock,
				{
					Query:    "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1",
					QueryArg: []driver.Value{int64(2)},
					Reply: helpers.MockReply{
						Columns: []string{"id", "name"},
						Rows: [][]driver.Value{
							{helpers.UserId, "User Name"},
						},
					},
				},
			},
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "error",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			helpers.PrepareDatabaseQueryMocks(test, sqlMock)
			db := helpers.InitGorm(dbMock)
			s := helpers.NewServer(db)

			c, recorder := helpers.PrepareContextFromTestCase(s, test)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				if assert.Equal(t, test.Expected.StatusCode, recorder.Code) {
					if recorder.Code == http.StatusOK {
						assertTokenResponse(t, recorder)
					}
				}
			}
		})
	}
}

func assertTokenResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var authResponse responses.LoginResponse
	_ = json.Unmarshal([]byte(recorder.Body.String()), &authResponse)

	assert.Equal(t, float64(helpers.UserId), getUserIdFromToken(authResponse.AccessToken))
	assert.Equal(t, float64(helpers.UserId), getUserIdFromToken(authResponse.RefreshToken))
}

func getUserIdFromToken(tokenToParse string) float64 {
	jwtToken, _ := jwt.Parse(tokenToParse, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		var hmacSampleSecret []byte
		return hmacSampleSecret, nil
	})
	claims, _ := jwtToken.Claims.(jwt.MapClaims)

	return claims["id"].(float64)
}
