package model

import (
	"time"
)

type SingleRating struct {
	GuestId string    `json:"guestId,omitempty"`
	Rating  int32     `json:"rating"`
	Date    time.Time `json:"date"`
}
