package dto

type SearchAccommodationDto struct {
	Amenities     []string `json:"amenities,omitempty"`
	EndDate       int64    `json:"endDate,omitempty"`
	StartDate     int64    `json:"startDate,omitempty"`
	Location      string   `json:"location,omitempty"`
	MaxPrice      int32    `json:"maxPrice,omitempty"`
	MinGuests     int32    `json:"minGuests,omitempty"`
	MinPrice      int32    `json:"minPrice,omitempty"`
	ProminentHost bool     `json:"prominentHost,omitempty"`
	MinRating     float32  `json:"minRating,omitempty"`
}
