package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniqueVisitor struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	IpAddress string             `json:"ipAddress" binding:"required" bson:"ipAddress"`
	Timestamp string             `json:"timestamp" binding:"required" bson:"timestamp"`
	Browser   string             `json:"browser" binding:"required" bson:"browser"`
}
