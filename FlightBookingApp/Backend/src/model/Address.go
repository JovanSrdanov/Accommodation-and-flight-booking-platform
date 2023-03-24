package model

type Address struct {
	Country      string `json:"country" binding:"required" bson:"country"`
	City         string `json:"city" binding:"required" bson:"city"`
	Street       string `json:"street" binding:"required" bson:"street"`
	StreetNumber string `json:"streetNumber" binding:"required" bson:"streetNumber"`
}
