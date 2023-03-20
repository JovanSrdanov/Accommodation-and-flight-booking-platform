package model

type Address struct {
	Country      string `json:"country" binding:"required"`
	City         string `json:"city" binding:"required"`
	Street       string `json:"street" binding:"required"`
	StreetNumber string `json:"streetNumber" binding:"required"`
}
