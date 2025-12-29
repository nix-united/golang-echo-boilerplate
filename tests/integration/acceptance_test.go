package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptance(t *testing.T) {
	registerRequest := requests.RegisterRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "example@email.com",
			Password: "some-password",
		},
		Name: "example-name",
	}

	rawRegisterRequest, err := json.Marshal(registerRequest)
	require.NoError(t, err)

	loginRequest := registerRequest.BasicAuth

	rawLoginRequest, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	createPostRequest := requests.CreatePostRequest{
		BasicPost: requests.BasicPost{
			Title:   "Title",
			Content: "Content",
		},
	}

	rawCreatePostRequest, err := json.Marshal(createPostRequest)
	require.NoError(t, err)

	var accessToken string

	t.Run("It should register an user", func(t *testing.T) {
		httpResponse, err := http.Post(
			applicationURL.JoinPath("/register").String(),
			"application/json",
			bytes.NewReader(rawRegisterRequest),
		)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusCreated, httpResponse.StatusCode)
	})

	t.Run("It should login", func(t *testing.T) {
		httpResponse, err := http.Post(
			applicationURL.JoinPath("/login").String(),
			"application/json",
			bytes.NewReader(rawLoginRequest),
		)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusOK, httpResponse.StatusCode)

		rawResponse, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		var loginResponse responses.LoginResponse
		err = json.Unmarshal(rawResponse, &loginResponse)
		require.NoError(t, err)

		require.NotEmpty(t, loginResponse.AccessToken)
		require.NotEmpty(t, loginResponse.Exp)
		require.NotEmpty(t, loginResponse.RefreshToken)

		accessToken = loginResponse.AccessToken
	})

	t.Run("It should create a post", func(t *testing.T) {
		httpRequest, err := http.NewRequest(
			http.MethodPost,
			applicationURL.JoinPath("/posts").String(),
			bytes.NewReader(rawCreatePostRequest),
		)
		require.NoError(t, err)

		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", "Bearer "+accessToken)

		httpResponse, err := http.DefaultClient.Do(httpRequest)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusCreated, httpResponse.StatusCode)
	})
}
