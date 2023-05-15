package communication

import "github.com/google/uuid"

type SimpleResponse struct {
	Message string `json:"message"`
}

func NewSimpleResponse(message string) *SimpleResponse {
	return &SimpleResponse{Message: message}
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

func NewErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{ErrorMessage: error}
}

type CreatedUserResponse struct {
	Username      string    `json:"username"`
	UserProfileId uuid.UUID `json:"userProfileId"`
}

func NewCreatedUserResponse(username string, userProfileId uuid.UUID) *CreatedUserResponse {
	return &CreatedUserResponse{Username: username, UserProfileId: userProfileId}
}
