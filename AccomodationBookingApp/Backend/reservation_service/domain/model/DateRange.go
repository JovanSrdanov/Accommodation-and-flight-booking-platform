package model

import (
	"time"
)

type DateRange struct {
	From time.Time `json:"from" binding:"required" bson:"from"`
	To   time.Time `json:"to" binding:"required" bson:"to"`
}

func (dateRange1 DateRange) IsInside(dateRange2 DateRange) bool {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if (dateRange1.From.After(dateRange2.From) || dateRange1.From.Equal(dateRange2.From)) && (dateRange1.To.Before(dateRange2.To) || dateRange1.To.Equal(dateRange2.To)) {
		return true
	}
	return false
}

func (dateRange1 DateRange) Overlaps(dateRange2 DateRange) bool {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if dateRange1.To.Before(dateRange2.From) || dateRange2.To.Before(dateRange1.From) {
		return false
	}
	return true
}

func (dateRange1 DateRange) IsStartFor(dateRange2 DateRange) bool {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if (dateRange1.From.Before(dateRange2.From) || dateRange1.From.Equal(dateRange2.From)) && (dateRange1.To.Before(dateRange2.To) || dateRange1.To.Equal(dateRange2.To)) {
		return true
	}
	return false
}

func (dateRange1 DateRange) IsEndFor(dateRange2 DateRange) bool {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if (dateRange1.To.After(dateRange2.To) || dateRange1.To.Equal(dateRange2.To)) && (dateRange1.From.After(dateRange2.From) || dateRange1.From.Equal(dateRange2.From)) {
		return true
	}
	return false
}

func (dateRange1 DateRange) Extends(dateRange2 DateRange) bool {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if dateRange1.From.Equal(dateRange2.To.AddDate(0, 0, 1).In(time.UTC)) {
		return true
	}
	return false
}

func (dateRange1 DateRange) DaysInCommon(dateRange2 DateRange) int32 {
	//NormalizeTime(&dateRange1)
	//NormalizeTime(&dateRange2)

	if !dateRange1.Overlaps(dateRange2) {
		return 0
	}

	// Find the start date and end date of the overlapping period
	startDate := dateRange1.From
	if dateRange2.From.After(startDate) {
		startDate = dateRange2.From
	}
	endDate := dateRange1.To
	if dateRange2.To.Before(endDate) {
		endDate = dateRange2.To
	}

	// Calculate the number of days between the start date and end date
	duration := endDate.Sub(startDate)
	days := int32(duration.Hours() / 24)

	return days
}

func NormalizeTime(dateRange *DateRange) {
	dateRange.To = dateRange.To.In(time.UTC)
	dateRange.From = dateRange.From.In(time.UTC)

	dateRange.To = time.Date(dateRange.To.Year(), dateRange.To.Month(), dateRange.To.Day(), 0, 0, 0, 0, dateRange.To.Location())
	dateRange.From = time.Date(dateRange.From.Year(), dateRange.From.Month(), dateRange.From.Day(), 0, 0, 0, 0, dateRange.From.Location())
}
