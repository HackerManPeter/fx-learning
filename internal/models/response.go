package models

type response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewErrorResponse(message string) response {
	return response{
		Status:  false,
		Message: message,
		Data:    nil,
	}
}

func NewSuccessResponse(data any, message string) response {
	return response{
		Status:  true,
		Message: message,
		Data:    data,
	}
}
