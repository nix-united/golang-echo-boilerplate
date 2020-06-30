package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/tests/helpers"
	"encoding/json"
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

type QueryMock struct {
	Query string
	Reply []map[string]interface{}
}

type TestCase struct{
	TestName  string
	Request   requests.RegisterRequest
	QueryMock *QueryMock
	Expected  ExpectedResponse
}

func TestWalk(t *testing.T) {
	cases := []TestCase {
		{
			"Register user success",
			requests.RegisterRequest{
				Name:     "name",
				Password: "password",
			},
			nil,
			ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "User successfully created",
			},
		},
		{
			"Register user with empty name",
			requests.RegisterRequest{
				Name:     "",
				Password: "password",
			},
			nil,
			ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
		{
			"Register user with too short password",
			requests.RegisterRequest{
				Name:     "name",
				Password: "passw",
			},
			nil,
			ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
		{
			"Register user with duplicated name",
			requests.RegisterRequest{
				Name:     "Duplicated Name",
				Password: "password",
			},
			&QueryMock{
				Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((name = Duplicated Name))`,
				Reply: []map[string]interface{}{{"id": 1, "name": "Duplicated Name", "password": "EncryptedPassword"}},
			},
			ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "User already exists",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			s := helpers.NewServer()

			requestJson, _ := json.Marshal(test.Request)
			request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(requestJson)))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			c := s.Echo.NewContext(request, recorder)

			if test.QueryMock != nil {
				mocket.Catcher.Reset().NewMock().WithQuery(test.QueryMock.Query).WithReply(test.QueryMock.Reply)
			}

			h := handlers.NewRegisterHandler(s)
			if assert.NoError(t, h.Register(c)) {
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
			}
		})
	}
}
