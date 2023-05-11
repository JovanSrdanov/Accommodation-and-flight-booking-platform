package model

type Amenity struct {
	Name string `json:"name" bson:"name"`
}

type Amenities []*Amenity
