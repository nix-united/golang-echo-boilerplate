package tests

import (
	"database/sql/driver"
	"echo-demo-project/requests"
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/tests/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWalkRegister(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		panic(err.Error())
	}

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
					Email:    "new-user@test.com",
					Password: "password",
				},
				Name: "name",
			},
			handlerFunc,
			[]*helpers.QueryMock{&helpers.SelectVersionMock,
				{
					Query:    "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL",
					QueryArg: []driver.Value{"new-user@test.com"},
					Reply: helpers.MockReply{
						Columns: []string{"id"},
					},
				},
			},
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
				Name: "",
			},
			handlerFunc,
			[]*helpers.QueryMock{
				&helpers.SelectVersionMock,
			},
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
				Name: "name",
			},
			handlerFunc,
			[]*helpers.QueryMock{
				&helpers.SelectVersionMock,
			},
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
				Name: "Another Name",
			},
			handlerFunc,
			[]*helpers.QueryMock{
				&helpers.SelectVersionMock,
				{
					Query:    "SELECT * FROM `users`  WHERE email = ? AND `users`.`deleted_at` IS NULL",
					QueryArg: []driver.Value{"duplicated@test.com"},
					Reply: helpers.MockReply{
						Columns: []string{"id", "email", "password"},
						Rows: [][]driver.Value{
							{helpers.UserId, "duplicated@test.com", "EncryptedPassword"},
						},
					},
				},
			},
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "User already exists",
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
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
			}
		})
	}
}
