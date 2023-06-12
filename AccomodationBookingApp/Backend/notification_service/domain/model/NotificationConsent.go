package model

import "github.com/google/uuid"

type NotificationConsent struct {
	UserProfileID            uuid.UUID `json:"userProfileID"  gorm:"primaryKey"`
	RequestMade              bool      `json:"requestMade"`
	ReservationCanceled      bool      `json:"ReservationCanceled"`
	HostRatingGiven          bool      `json:"HostRatingGiven"`
	AccommodationRatingGiven bool      `json:"AccommodationRatingGiven"`
	ProminentHost            bool      `json:"AccommodationRatingGiven"`
	HostResponded            bool      `json:"HostResponded"`
}

func (NotificationConsent) TableName() string {
	return "notification_consent"
}