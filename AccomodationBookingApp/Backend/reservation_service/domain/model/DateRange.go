package model

import (
	"time"
)

type DateRange struct {
	From time.Time `json:"from" binding:"required" bson:"from"`
	To   time.Time `json:"to" binding:"required" bson:"to"`
}

func (dateRange1 DateRange) IsInside(dateRange2 DateRange) bool {
	NormalizeTime(dateRange1)
	NormalizeTime(dateRange2)

	if (dateRange1.From.After(dateRange2.From) || dateRange1.From.Equal(dateRange2.From)) && (dateRange1.To.Before(dateRange2.To) || dateRange1.To.Equal(dateRange2.To)) {
		return true
	}
	return false
}

func (dateRange1 DateRange) Overlaps(dateRange2 DateRange) bool {
	NormalizeTime(dateRange1)
	NormalizeTime(dateRange2)

	if dateRange1.To.Before(dateRange2.From) || dateRange2.To.Before(dateRange1.From) {
		return false
	}
	return true
}

func NormalizeTime(dateRange DateRange) {
	dateRange.To = time.Date(dateRange.To.Year(), dateRange.To.Month(), dateRange.To.Day(), 0, 0, 0, 0, dateRange.To.Location())
	dateRange.From = time.Date(dateRange.From.Year(), dateRange.From.Month(), dateRange.From.Day(), 0, 0, 0, 0, dateRange.From.Location())

	dateRange.To = dateRange.To.In(time.UTC)
	dateRange.From = dateRange.From.In(time.UTC)
}
