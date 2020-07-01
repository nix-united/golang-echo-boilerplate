package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/models"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"echo-demo-project/tests/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalkAuth(t *testing.T) {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((name = User Name))`,
		Reply: helpers.MockReply{{"id": helpers.UserId, "name": "User Name", "password": encryptedPassword}},
	}

	cases := []helpers.TestCase {
		{
			"Auth success",
			requests.LoginRequest{
				Name:     "User Name",
				Password: "password",
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "",
			},
		},
		{
			"Login attempt with incorrect password",
			requests.LoginRequest{
				Name:     "User Name",
				Password: "incorrectPassword",
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
		{
			"Login attempt as non-existent user",
			requests.LoginRequest{
				Name:     "User Not Exists",
				Password: "password",
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test, "/login")
			h := handlers.NewAuthHandler(s)

			if assert.NoError(t, h.Login(c)) {
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
	tokenService := services.NewTokenService()

	validUser := models.User{Name: "User Name"}
	validUser.ID = helpers.UserId
	validToken, _ := tokenService.CreateRefreshToken(&validUser)

	notExistUser := models.User{Name: "User Not Exists"}
	notExistUser.ID = helpers.UserId + 1
	notExistToken, _ := tokenService.CreateRefreshToken(&notExistUser)

	invalidToken := validToken[1:len(validToken)-1]

	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 1))`,
		Reply: helpers.MockReply{{"id": helpers.UserId, "name": "User Name"}},
	}

	cases := []helpers.TestCase {
		{
			"Refresh success",
			requests.RefreshRequest{
				Token: validToken,
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "",
			},
		},
		{
			"Refresh token of non-existent user",
			requests.RefreshRequest{
				Token: notExistToken,
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "User not found",
			},
		},
		{
			"Refresh invalid token",
			requests.RefreshRequest{
				Token: invalidToken,
			},
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "error",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test, "/refresh")
			h := handlers.NewAuthHandler(s)

			if assert.NoError(t, h.RefreshToken(c)) {
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
	token, _ := jwt.Parse(tokenToParse, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		var hmacSampleSecret []byte
		return hmacSampleSecret, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["id"].(float64)
}