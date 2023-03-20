package dto

type SimpleResponse struct {
	Message string `json:"message"`
}

func NewSimpleResponse(message string) *SimpleResponse {
	return &SimpleResponse{
		Message: message,
	}
}
