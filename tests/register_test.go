package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/tests/helpers"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Response struct {
	StatusCode int
	BodyPart   string
}

func TestWalk(t *testing.T) {
	cases := []struct{
		Name	 string
		Request  requests.RegisterRequest
		Expected Response
	} {
		{
			"Register user success",
			requests.RegisterRequest{
				Name:     "name",
				Password: "password",
			},
			Response{
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
			Response{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
		{
			"Register user with too short password",
			requests.RegisterRequest{
				Name:     "name",
				Password:  "passw",
			},
			Response{
				StatusCode: 400,
				BodyPart:   "error",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			s := helpers.NewServer()

			requestJson, _ := json.Marshal(test.Request)
			request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(requestJson)))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			c := s.Echo.NewContext(request, recorder)

			h := handlers.NewRegisterHandler(s)
			if assert.NoError(t, h.Register(c)) {
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
			}
		})
	}
}
