package dto

type HostRating struct {
	AvgRating float32             `json:"avgRating"`
	HostID    string              `json:"hostId"`
	Ratings   []*HostRatingRating `json:"ratings"`
}

type HostRatingRating struct {
	GuestID string `json:"guestId"`
	Date    string `json:"Date"`
	Rating  int32  `json:"rating"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
