package model

import "time"

type DateRange struct {
	From time.Time `json:"from" binding:"required" bson:"from"`
	To   time.Time `json:"to" binding:"required" bson:"to"`
}
