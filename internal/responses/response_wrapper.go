package responses

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

func NewErrorResponse(message string, code int) ErrorResponse {
	return ErrorResponse{
		Code:  code,
		Error: message,
	}
}
