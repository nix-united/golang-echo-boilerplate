package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func newPostHandler(t *testing.T) (*handlers.PostHandlers, *MockpostService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	postService := NewMockpostService(ctrl)
	postHandler := handlers.NewPostHandlers(postService)

	return postHandler, postService
}

func TestPostHandler_CreatePost(t *testing.T) {
	const userID = 200

	authClaims := &jwt.Token{Claims: &token.JwtCustomClaims{
		ID:   userID,
		Name: "user_name",
	}}

	request := requests.CreatePostRequest{
		BasicPost: requests.BasicPost{
			Title:   "post-title",
			Content: "post-content",
		},
	}

	invalidRequest := requests.CreatePostRequest{
		BasicPost: requests.BasicPost{
			Content: "post-content",
		},
	}

	wantPost := &models.Post{
		Title:   "post-title",
		Content: "post-content",
		UserID:  userID,
	}

	testCases := map[string]struct {
		setExpectations func(postService *MockpostService)
		request         any
		wantStatus      int
		wantResponse    any
	}{
		"It should respond with 400 status code if request is invalid": {
			setExpectations: func(postService *MockpostService) {},
			request:         invalidRequest,
			wantStatus:      http.StatusBadRequest,
			wantResponse: responses.Error{
				Code:  http.StatusBadRequest,
				Error: "Required fields are empty",
			},
		},
		"It should create a post": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, gotPost *models.Post) error {
						assert.Equal(t, wantPost, gotPost)

						return nil
					})
			},
			request:    request,
			wantStatus: http.StatusCreated,
			wantResponse: responses.Data{
				Code:    http.StatusCreated,
				Message: "Post successfully created",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			postHandler, postService := newPostHandler(t)

			testCase.setExpectations(postService)

			rawRequest, err := json.Marshal(testCase.request)
			require.NoError(t, err)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/posts",
				bytes.NewBuffer(rawRequest),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			c.Set("user", authClaims)

			err = postHandler.CreatePost(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}

func TestPostHandler_GetPosts(t *testing.T) {
	postHandler, postService := newPostHandler(t)

	posts := []models.Post{{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "post-title",
		Content: "post-content",
		UserID:  200,
		User: models.User{
			Model: gorm.Model{
				ID: 200,
			},
			Email:    "example@email.com",
			Name:     "example-name",
			Password: "password",
		},
	}}

	wantResponse := []responses.PostResponse{{
		ID:       100,
		Title:    "post-title",
		Content:  "post-content",
		Username: "example-name",
	}}

	postService.
		EXPECT().
		GetPosts(gomock.Any()).
		Return(posts, nil)

	request := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"/posts",
		http.NoBody,
	)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	recorder := httptest.NewRecorder()
	c := echo.New().NewContext(request, recorder)

	err := postHandler.GetPosts(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	rawWantResponse, err := json.Marshal(wantResponse)
	require.NoError(t, err)

	assert.JSONEq(t, string(rawWantResponse), recorder.Body.String())
}

func TestPostHandler_UpdatePost(t *testing.T) {
	const postOwnerID = 200

	authClaims := &jwt.Token{Claims: &token.JwtCustomClaims{
		ID:   postOwnerID,
		Name: "user_name",
	}}

	post := models.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "post-title",
		Content: "post-content",
		UserID:  postOwnerID,
	}

	postWithDifferentUser := models.Post{
		Model: gorm.Model{
			ID: 100,
		},
		UserID: 201,
	}

	request := requests.UpdatePostRequest{
		BasicPost: requests.BasicPost{
			Title:   "new-title",
			Content: "new-content",
		},
	}

	rawRequest, err := json.Marshal(request)
	require.NoError(t, err)

	testCases := map[string]struct {
		setExpectations func(postService *MockpostService)
		wantStatus      int
		wantResponse    any
	}{
		"It should return a 404 status code when post not found": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(models.Post{}, models.ErrPostNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantResponse: responses.Error{
				Code:  http.StatusNotFound,
				Error: "Post not found",
			},
		},
		"It should return a 403 status code when user tries to update not their post": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(postWithDifferentUser, nil)
			},
			wantStatus: http.StatusForbidden,
			wantResponse: responses.Error{
				Code:  http.StatusForbidden,
				Error: "Forbidden",
			},
		},
		"It should update post": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(post, nil)

				postService.
					EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						_ context.Context,
						gotPost *models.Post,
						gotRequest requests.UpdatePostRequest,
					) error {
						assert.Equal(t, &post, gotPost)
						assert.Equal(t, request, gotRequest)

						return nil
					})
			},
			wantStatus: http.StatusOK,
			wantResponse: responses.Data{
				Code:    http.StatusOK,
				Message: "Post successfully updated",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			postHandler, postService := newPostHandler(t)

			testCase.setExpectations(postService)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPut,
				fmt.Sprintf("/posts/%d", post.ID),
				bytes.NewBuffer(rawRequest),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			c.SetPath("/posts/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatUint(uint64(post.ID), 10))
			c.Set("user", authClaims)

			err := postHandler.UpdatePost(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}

func TestPostHandler_DeletePost(t *testing.T) {
	const postOwnerID = 200

	authClaims := &jwt.Token{Claims: &token.JwtCustomClaims{
		ID:   postOwnerID,
		Name: "user_name",
	}}

	post := models.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "post-title",
		Content: "post-content",
		UserID:  postOwnerID,
	}

	postWithDifferentUser := models.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "post-title",
		Content: "post-content",
		UserID:  201,
	}

	testCases := map[string]struct {
		setExpectations func(postService *MockpostService)
		wantStatus      int
		wantResponse    any
	}{
		"It should return 404 status code when post not found": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(models.Post{}, models.ErrPostNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantResponse: responses.Error{
				Code:  http.StatusNotFound,
				Error: "Post not found",
			},
		},
		"It should return 403 status code when user tries to delete not their post": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(postWithDifferentUser, nil)
			},
			wantStatus: http.StatusForbidden,
			wantResponse: responses.Error{
				Code:  http.StatusForbidden,
				Error: "Forbidden",
			},
		},
		"It should delete post": {
			setExpectations: func(postService *MockpostService) {
				postService.
					EXPECT().
					GetPost(gomock.Any(), post.ID).
					Return(post, nil)

				postService.
					EXPECT().
					Delete(gomock.Any(), &post).
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
			wantResponse: responses.Data{
				Code:    http.StatusNoContent,
				Message: "Post deleted successfully",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			postHandler, postService := newPostHandler(t)

			testCase.setExpectations(postService)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodDelete,
				fmt.Sprintf("/posts/%d", post.ID),
				http.NoBody,
			)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			c.SetPath("/posts/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatUint(uint64(post.ID), 10))
			c.Set("user", authClaims)

			err := postHandler.DeletePost(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}
