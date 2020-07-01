package tests

import (
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/services"
	"echo-demo-project/tests/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWalkPostsCrud(t *testing.T) {
	requestCreate := helpers.Request{
		Method: http.MethodPost,
		Url:    "/posts",
	}
	requestGet := helpers.Request{
		Method: http.MethodGet,
		Url:    "/posts",
	}
	requestUpdate := helpers.Request{
		Method: http.MethodPut,
		Url:    "/posts/",
	}
	requestDelete := helpers.Request{
		Method: http.MethodDelete,
		Url:    "/posts/",
	}
	handlerFuncCreate := func(s *server.Server, c echo.Context) error {
		return handlers.NewPostHandlers(s).CreatePost(c)
	}
	handlerFuncGet := func(s *server.Server, c echo.Context) error {
		return handlers.NewPostHandlers(s).GetPosts(c)
	}
	handlerFuncUpdate := func(s *server.Server, c echo.Context) error {
		return handlers.NewPostHandlers(s).UpdatePost(c)
	}
	handlerFuncDelete := func(s *server.Server, c echo.Context) error {
		return handlers.NewPostHandlers(s).DeletePost(c)
	}

	claims := &services.JwtCustomClaims{
		Name: "user",
		ID:   helpers.UserId,
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "posts"  WHERE `,
		Reply: helpers.MockReply{{"id": 1, "title": "title", "content": "content", "username": "Username"}},
	}

	cases := []helpers.TestCase {
		{
			"Create post success",
			requestCreate,
			requests.CreatePostRequest{
				Title:   "title",
				Content: "content",
			},
			handlerFuncCreate,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Post successfully created",
			},
		},
		{
			"Create post with empty title",
			requestCreate,
			requests.CreatePostRequest{
				Title:   "",
				Content: "content",
			},
			handlerFuncCreate,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Required fields are empty",
			},
		},
		{
			"Get posts success",
			requestGet,
			"",
			handlerFuncGet,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "[{\"title\":\"title\",\"content\":\"content\",\"username\":\"\",\"id\":1}]",
			},
		},
		{
			"Update post success",
			requestUpdate,
			requests.UpdatePostRequest{
				Title:   "new title",
				Content: "new content",
			},
			handlerFuncUpdate,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Post successfully updated",
			},
		},
		{
			"Update post with empty title",
			requestUpdate,
			requests.UpdatePostRequest{
				Title:   "",
				Content: "new content",
			},
			handlerFuncUpdate,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Required fields are empty",
			},
		},
		{
			"Delete post success",
			requestDelete,
			"",
			handlerFuncDelete,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Post deleted successfully",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test)
			c.Set("user", validToken)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
			}
		})
	}
}
