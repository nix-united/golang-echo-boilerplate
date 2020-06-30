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
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
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

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			s := helpers.NewServer()

			requestJson, _ := json.Marshal(test.Request)
			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(requestJson)))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			c := s.Echo.NewContext(request, recorder)

			if test.QueryMock != nil {
				mocket.Catcher.Reset().NewMock().WithQuery(test.QueryMock.Query).WithReply(test.QueryMock.Reply)
			}

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

	token, _ := jwt.Parse(authResponse.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		var hmacSampleSecret []byte
		return hmacSampleSecret, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	assert.Equal(t, float64(userId), claims["id"])
}
