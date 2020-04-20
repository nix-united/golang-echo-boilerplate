package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/services"
	"echo-demo-project/tests/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreatePost(t *testing.T)  {
	s := helpers.NewServer()

	f := make(url.Values)
	f.Set("title", "title")
	f.Set("content", "content")

	req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	claims := &services.JwtCustomClaims{
		Name: "user",
		ID:    1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	c.Set("user", token)

	h := handlers.NewPostHandlers(s)
	if assert.NoError(t, h.CreatePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Post successfully create\"", strings.TrimSpace(rec.Body.String()))
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
		assert.Equal(t, "\"Post delete successfully\"", strings.TrimSpace(rec.Body.String()))
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
		assert.Equal(t, "[{\"Title\":\"title\",\"Content\":\"content\",\"Username\":\"\",\"ID\":1}]", strings.TrimSpace(rec.Body.String()))
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
		assert.Equal(t, "\"Post successfully update\"", strings.TrimSpace(rec.Body.String()))
	}
}