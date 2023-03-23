package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatedResponse struct {
	Id primitive.ObjectID `json:"id"`
}

func NewCreatedResponse(id primitive.ObjectID) *CreatedResponse {
	return &CreatedResponse{
		Id: id,
	}
}
