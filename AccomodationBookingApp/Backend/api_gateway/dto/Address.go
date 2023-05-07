package dto

type Address struct {
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	Street       string `json:"street,omitempty"`
	StreetNumber string `json:"streetNumber,omitempty"`
}
