package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/server/responses"
	"echo-demo-project/tests/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)


func TestAuth(t *testing.T) {
	s := helpers.NewServer()

	f := make(url.Values)
	f.Set("name", "name")
	f.Set("password", "password")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	commonReply := []map[string]interface{}{
		{"id": 1, "name": "name", "password": string(encryptedPassword)},
	}
	mocket.Catcher.Reset().NewMock().WithArgs("name").WithReply(commonReply)

	var authResponse responses.LoginResponse

	h := handlers.NewAuthHandler(s)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		_ = json.Unmarshal([]byte(rec.Body.String()), &authResponse)

		token, _ := jwt.Parse(authResponse.AccessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			}
			var hmacSampleSecret []byte
			return hmacSampleSecret, nil
		})
		claims, _ := token.Claims.(jwt.MapClaims)

		assert.Equal(t, float64(1), claims["id"])
	}
}
