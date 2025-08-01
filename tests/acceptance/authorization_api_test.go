package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationAPI(t *testing.T) {
	const (
		email    = "example@email.com"
		password = "some-password"
	)

	t.Run("It should register an user", func(t *testing.T) {
		registerRequest := requests.RegisterRequest{
			BasicAuth: requests.BasicAuth{
				Email:    email,
				Password: password,
			},
			Name: "some-name",
		}

		requestBody := new(bytes.Buffer)
		err := json.NewEncoder(requestBody).Encode(registerRequest)
		require.NoError(t, err)

		response, err := http.Post(
			applicationURL.JoinPath("/register").String(),
			"application/json",
			requestBody,
		)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, response.Body.Close())
		}()

		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("It should login an user", func(t *testing.T) {
		registerRequest := requests.LoginRequest{
			BasicAuth: requests.BasicAuth{
				Email:    email,
				Password: password,
			},
		}

		requestBody := new(bytes.Buffer)
		err := json.NewEncoder(requestBody).Encode(registerRequest)
		require.NoError(t, err)

		response, err := http.Post(
			applicationURL.JoinPath("/login").String(),
			"application/json",
			requestBody,
		)
		defer func() {
			assert.NoError(t, response.Body.Close())
		}()

		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)

		var loginResponse responses.LoginResponse
		err = json.NewDecoder(response.Body).Decode(&loginResponse)
		require.NoError(t, err)

		assert.NotEmpty(t, loginResponse.AccessToken)
	})
}
