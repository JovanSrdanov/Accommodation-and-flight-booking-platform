package model

type RatingResponse struct {
	AccommodationId string  `json:"accommodationId"`
	Rating          float32 `json:"rating"`
}
