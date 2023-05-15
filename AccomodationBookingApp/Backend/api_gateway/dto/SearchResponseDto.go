package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Address   Address            `json:"address,omitempty"`
	MinGuests int32              `json:"minGuests,omitempty"`
	MaxGuests int32              `json:"maxGuests,omitempty"`
	Amenities []string           `json:"amenities,omitempty"`
	Images    []string           `json:"images,omitempty"`
	HostId    string             `json:"hostId,omitempty"`
	Price     int32              `json:"price,omitempty"`
}

type SearchResponseDto []*Accommodation
