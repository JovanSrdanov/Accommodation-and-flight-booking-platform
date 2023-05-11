package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceWithDate struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	DateRange        DateRange          `json:"dateRange" binding:"required" bson:"dateRange"`
	Price            int32              `json:"price" binding:"required" bson:"price"`
	IsPricePerPerson bool               `json:"isPricePerPerson" binding:"required" bson:"isPricePerPerson"`
}
