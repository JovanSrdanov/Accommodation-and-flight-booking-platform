package model

type RatingResponse struct {
	AccommodationId string          `json:"accommodationId"`
	AvgRating       float32         `json:"avgRating"`
	Ratings         []*SingleRating `json:"ratings"`
}

type SimpleRatingResponse struct {
	AccommodationId string  `json:"accommodationId"`
	AvgRating       float32 `json:"avgRating"`
}
