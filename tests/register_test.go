package tests

import (
	"echo-demo-project/server/handlers"
	"echo-demo-project/tests/helpers"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)


func TestCreateUser(t *testing.T)  {
	s := helpers.NewServer()

	f := make(url.Values)
	f.Set("name", "name")
	f.Set("password", "password")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	h := handlers.NewRegisterHandler(s)
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"User successfully create\"", strings.TrimSpace(rec.Body.String()))
	}
}
