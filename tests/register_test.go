package tests

import (
	"echo-demo-project/requests"
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/tests/helpers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWalkRegister(t *testing.T) {
	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/register",
	}
	handlerFunc := func(s *server.Server, c echo.Context) error {
		return handlers.NewRegisterHandler(s).Register(c)
	}

	cases := []helpers.TestCase{
		{
			"Register user success",
			request,
			requests.RegisterRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "password",
				},
				Name: "name",
			},
			handlerFunc,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "User successfully created",
			},
		},
		{
			"Register user with empty name",
			request,
			requests.RegisterRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "password",
				},
				Name:     "",
			},
			handlerFunc,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
		{
			"Register user with too short password",
			request,
			requests.RegisterRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "passw",
				},
				Name:     "name",
			},
			handlerFunc,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
		{
			"Register user with duplicated email",
			request,
			requests.RegisterRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "duplicated@test.com",
					Password: "password",
				},
				Name:     "Another Name",
			},
			handlerFunc,
			&helpers.QueryMock{
				Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((email = duplicated@test.com))`,
				Reply: helpers.MockReply{{"id": 1, "email": "duplicated@test.com", "password": "EncryptedPassword"}},
			},
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "User already exists",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
			}
		})
	}
}
