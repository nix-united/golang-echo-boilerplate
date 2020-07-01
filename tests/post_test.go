package tests

import (
	"echo-demo-project/server"
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/services"
	"echo-demo-project/tests/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWalkCreatePost(t *testing.T) {
	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/posts",
	}
	handlerFunc := func(s *server.Server, c echo.Context) error {
		return handlers.NewPostHandlers(s).CreatePost(c)
	}

	claims := &services.JwtCustomClaims{
		Name: "user",
		ID:   helpers.UserId,
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cases := []helpers.TestCase {
		{
			"Create post success",
			request,
			requests.CreatePostRequest{
				Title:   "title",
				Content: "content",
			},
			handlerFunc,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Post successfully created",
			},
		},
		{
			"Create post with empty title",
			request,
			requests.CreatePostRequest{
				Title:   "",
				Content: "content",
			},
			handlerFunc,
			nil,
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Required fields are empty",
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

func TestDeletePost(t *testing.T)  {
	s := helpers.NewServer()

	req := httptest.NewRequest(http.MethodDelete, "/posts", nil)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	c.Set("id", 1)
	commonReply := []map[string]interface{}{{"id": 1, "title": "title", "content": "content"}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "posts"  WHERE `).WithReply(commonReply)

	h := handlers.NewPostHandlers(s)
	if assert.NoError(t, h.DeletePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Post deleted successfully\"", strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetPosts(t *testing.T)  {
	s := helpers.NewServer()

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	commonReply := []map[string]interface{}{{"id": 1, "title": "title", "content": "content", "username": "Username"}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "posts"  WHERE `).WithReply(commonReply)

	h := handlers.NewPostHandlers(s)
	if assert.NoError(t, h.GetPosts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[{\"title\":\"title\",\"content\":\"content\",\"username\":\"\",\"id\":1}]", strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdatePost(t *testing.T)  {
	s := helpers.NewServer()
	f := make(url.Values)
	f.Set("title", "title")
	f.Set("content", "content")

	req := httptest.NewRequest(http.MethodPut, "/posts", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	commonReply := []map[string]interface{}{{"id": 1, "title": "title", "content": "content", "username": "Username"}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "posts"  WHERE `).WithReply(commonReply)

	h := handlers.NewPostHandlers(s)
	if assert.NoError(t, h.UpdatePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Post successfully updated\"", strings.TrimSpace(rec.Body.String()))
	}
}