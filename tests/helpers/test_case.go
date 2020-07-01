package helpers

import (
	"echo-demo-project/server"
	"encoding/json"
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"net/http"
	"net/http/httptest"
	"strings"
)

type MockReply []map[string]interface{}

type QueryMock struct {
	Query string
	Reply MockReply
}

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

type TestCase struct{
	TestName  string
	Request   interface{}
	QueryMock *QueryMock
	Expected  ExpectedResponse
}

func PrepareContextFromTestCase(s *server.Server, test TestCase, requestTarget string) (c echo.Context, recorder *httptest.ResponseRecorder) {
	requestJson, _ := json.Marshal(test.Request)
	request := httptest.NewRequest(http.MethodPost, requestTarget, strings.NewReader(string(requestJson)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder = httptest.NewRecorder()
	c = s.Echo.NewContext(request, recorder)

	if test.QueryMock != nil {
		mocket.Catcher.Reset().NewMock().WithQuery(test.QueryMock.Query).WithReply(test.QueryMock.Reply)
	}

	return
}