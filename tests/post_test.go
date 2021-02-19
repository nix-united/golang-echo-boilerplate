package tests

import (
	"echo-demo-project/models"
	"echo-demo-project/requests"
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/tests/helpers"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const postId = "1"
const postIdNotExists = "2"

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
		Url:    "/posts/" + postId,
		PathParam: &helpers.PathParam{
			Name:  "id",
			Value: postId,
		},
	}
	requestDelete := helpers.Request{
		Method: http.MethodDelete,
		Url:    "/posts/" + postId,
		PathParam: &helpers.PathParam{
			Name:  "id",
			Value: postId,
		},
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

	validUser := models.User{Email: "name@test.com"}
	validUser.ID = helpers.UserId

	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "posts"  WHERE "posts"."deleted_at" IS NULL AND ((id = 1 ))`,
		Reply: helpers.MockReply{{"id": 1, "title": "title", "content": "content", "username": "Username"}},
	}

	cases := []helpers.TestCase {
		{
			"Create post success",
			requestCreate,
			requests.CreatePostRequest{
				BasicPost: requests.BasicPost{
					Title:   "title",
					Content: "content",
				},
			},
			handlerFuncCreate,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "Post successfully created",
			},
		},
		{
			"Create post with empty title",
			requestCreate,
			requests.CreatePostRequest{
				BasicPost: requests.BasicPost{
					Title:   "",
					Content: "content",
				},
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
			&helpers.QueryMock{
				Query: `SELECT * FROM "posts"  WHERE `,
				Reply: helpers.MockReply{{"id": 1, "title": "title", "content": "content", "username": "Username"}},
			},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "[{\"title\":\"title\",\"content\":\"content\",\"username\":\"\",\"id\":1}]",
			},
		},
		{
			"Update post success",
			requestUpdate,
			requests.UpdatePostRequest{
				BasicPost: requests.BasicPost{
					Title:   "new title",
					Content: "new content",
				},
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
				BasicPost: requests.BasicPost{
					Title:   "",
					Content: "new content",
				},
			},
			handlerFuncUpdate,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Required fields are empty",
			},
		},
		{
			"Update non-existent post",
			helpers.Request{
				Method: http.MethodPut,
				Url:    "/posts/" + postIdNotExists,
				PathParam: &helpers.PathParam{
					Name:  "id",
					Value: postIdNotExists,
				},
			},
			requests.UpdatePostRequest{
				BasicPost: requests.BasicPost{
					Title:   "new title",
					Content: "new content",
				},
			},
			handlerFuncUpdate,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "Post not found",
			},
		},
		{
			"Delete post success",
			requestDelete,
			"",
			handlerFuncDelete,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 204,
				BodyPart:   "Post deleted successfully",
			},
		},
		{
			"Delete non-existent post",
			helpers.Request{
				Method: http.MethodDelete,
				Url:    "/posts/" + postIdNotExists,
				PathParam: &helpers.PathParam{
					Name:  "id",
					Value: postIdNotExists,
				},
			},
			"",
			handlerFuncDelete,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "Post not found",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test)
			c.Set("currentUser", &validUser)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
			}
		})
	}
}
