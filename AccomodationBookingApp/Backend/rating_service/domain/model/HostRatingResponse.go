package model

type HostRatingResponse struct {
	HostId    string          `json:"hostId"`
	AvgRating float32         `json:"avgRating"`
	Ratings   []*SingleRating `json:"ratings"`
}
