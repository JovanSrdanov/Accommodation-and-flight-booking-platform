package model

type Address struct {
	Country      string `json:"country""`
	City         string `json:"city"`
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
}
