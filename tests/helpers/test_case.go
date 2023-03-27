package helpers

import (
	"database/sql/driver"
	"echo-demo-project/server"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"regexp"
	"strings"
)

const UserId = 1

type TestCase struct {
	TestName    string
	Request     Request
	RequestBody interface{}
	HandlerFunc func(s *server.Server, c echo.Context) error
	QueryMocks  []*QueryMock
	Expected    ExpectedResponse
}

type PathParam struct {
	Name  string
	Value string
}

type Request struct {
	Method    string
	Url       string
	PathParam *PathParam
}

type MockReply struct {
	Columns []string
	Rows    [][]driver.Value
}

type QueryMock struct {
	Query    string
	QueryArg []driver.Value
	Reply    MockReply
}

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

var SelectVersionMock = QueryMock{
	Query: "SELECT VERSION()",
	Reply: MockReply{
		Columns: []string{"VERSION()"},
		Rows: [][]driver.Value{
			{"8.0.32"},
		},
	},
}

func PrepareContextFromTestCase(s *server.Server, test TestCase) (c echo.Context, recorder *httptest.ResponseRecorder) {
	requestJson, _ := json.Marshal(test.RequestBody)
	request := httptest.NewRequest(test.Request.Method, test.Request.Url, strings.NewReader(string(requestJson)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder = httptest.NewRecorder()
	c = s.Echo.NewContext(request, recorder)

	if test.Request.PathParam != nil {
		c.SetParamNames(test.Request.PathParam.Name)
		c.SetParamValues(test.Request.PathParam.Value)
	}

	return
}

func PrepareDatabaseQueryMocks(test TestCase, mock sqlmock.Sqlmock) {
	for _, queryMock := range test.QueryMocks {
		query := mock.ExpectQuery(regexp.QuoteMeta(queryMock.Query)).
			WillReturnRows(mock.NewRows(queryMock.Reply.Columns))

		if len(queryMock.QueryArg) != 0 {
			query.WithArgs(queryMock.QueryArg...)
		}

		rows := sqlmock.NewRows(queryMock.Reply.Columns)
		for _, row := range queryMock.Reply.Rows {
			rows.AddRow(row...)
		}
		query.WillReturnRows(rows)
	}
}
