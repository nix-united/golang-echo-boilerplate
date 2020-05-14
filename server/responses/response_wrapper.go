package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Code  int `json:"code"`
	Error string `json:"error"`
}

func Response(c echo.Context, statusCode int, data interface{}) error {
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return c.JSON(statusCode, data)
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return Response(c, http.StatusOK, data)
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Error{
		Code:  statusCode,
		Error: message,
	})
}
