package model

type Address struct {
	Country      string `json:"country" binding:"required" bson:"country" example:"Serbia"`
	City         string `json:"city" binding:"required" bson:"city" example:"Novi Sad"`
	Street       string `json:"street" binding:"required" bson:"street" example:"Rumenacka"`
	StreetNumber string `json:"streetNumber" binding:"required" bson:"streetNumber" example:"21a"`
}
