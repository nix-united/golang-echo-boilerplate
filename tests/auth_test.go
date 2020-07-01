package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
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

const userId = 1

func TestWalkAuth(t *testing.T) {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((name = User Name))`,
		Reply: helpers.MockReply{{"id": userId, "name": "User Name", "password": encryptedPassword}},
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

func assertTokenResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var authResponse responses.LoginResponse
	_ = json.Unmarshal([]byte(recorder.Body.String()), &authResponse)

	assert.Equal(t, float64(userId), getUserIdFromToken(authResponse.AccessToken))
	assert.Equal(t, float64(userId), getUserIdFromToken(authResponse.RefreshToken))
}

func getUserIdFromToken(tokenToParse string) interface{} {
	token, _ := jwt.Parse(tokenToParse, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		var hmacSampleSecret []byte
		return hmacSampleSecret, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["id"]
}