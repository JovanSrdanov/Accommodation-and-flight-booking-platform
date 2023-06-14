package dto

type AccommodationRating struct {
	AvgRating       float32                      `json:"avgRating"`
	AccommodationID string                       `json:"accommodationId"`
	Ratings         []*AccommodationRatingRating `json:"ratings"`
}

type AccommodationRatingRating struct {
	GuestID string `json:"guestId"`
	Date    string `json:"Date"`
	Rating  int32  `json:"rating"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
