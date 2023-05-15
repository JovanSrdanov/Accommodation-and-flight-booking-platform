package model

type SearchDto struct {
	HostId    string   `json:"hostId" binding:"required" bson:"hostId"`
	MinGuests int32    `json:"minGuests" binding:"required" bson:"minGuests"`
	Amenities []string `json:"amenities" binding:"required" bson:"amenities"`
	Location  string   `json:"location" binding:"required" bson:"location"`
}
