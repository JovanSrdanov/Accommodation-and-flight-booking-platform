package model

import (
	"time"
)

type RateHostDto struct {
	HostId  string    `json:"hostId"`
	GuestId string    `json:"guestId,omitempty"`
	Rating  int32     `json:"rating"`
	Date    time.Time `json:"date"`
}
