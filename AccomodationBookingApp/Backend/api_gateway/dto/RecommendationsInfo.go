package dto

type Recommendation struct {
	AccommodationID string   `json:"accommodationId"`
	Rating          float32  `json:"rating"`
	Name            string   `json:"name,omitempty"`
	Address         Address  `json:"address,omitempty"`
	MinGuests       int32    `json:"minGuests,omitempty"`
	MaxGuests       int32    `json:"maxGuests,omitempty"`
	Amenities       []string `json:"amenities,omitempty"`
	Images          []string `json:"images,omitempty"`
	HostId          string   `json:"hostId,omitempty"`
}
