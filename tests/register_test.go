package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/tests/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalkRegister(t *testing.T) {
	cases := []helpers.TestCase {
		{
			"Register user success",
			requests.RegisterRequest{
				Name:     "name",
				Password: "password",
			},
			nil,
			helpers.ExpectedResponse{
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
			helpers.ExpectedResponse{
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
			helpers.ExpectedResponse{
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
			&helpers.QueryMock{
				Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((name = Duplicated Name))`,
				Reply: helpers.MockReply{{"id": 1, "name": "Duplicated Name", "password": "EncryptedPassword"}},
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
			c, recorder := helpers.PrepareContextFromTestCase(s, test, "/register")
			h := handlers.NewRegisterHandler(s)

			if assert.NoError(t, h.Register(c)) {
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
			}
		})
	}
}
