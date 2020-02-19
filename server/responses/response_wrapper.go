package responses

import (
	"github.com/labstack/echo"
	"net/http"
)

func Response(c echo.Context, statusCode int, data interface{}) error {
	//context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	//context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return c.JSON(statusCode, data)
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return Response(c, http.StatusOK, data)
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, struct {
		Code  int
		Error string
	}{Code: statusCode, Error: message})
}
