package startup

import "os"

type Configuration struct {
	Port                               string
	DBName                             string
	DBPort                             string
	ReservationServiceHost             string
	ReservationServicePort             string
	NatsHost                           string
	NatsPort                           string
	NatsUser                           string
	NatsPass                           string
	DeleteUserCommandSubject           string
	DeleteUserReplySubject             string
	AccommodationEventDbName           string
	AccommodationEventDbPort           string
	AccommodationEventInnerDbName      string
	AccommodationEventDbCollectionName string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                               os.Getenv("ACCOMMODATION_SERVICE_PORT"),
		DBName:                             os.Getenv("ACCOMMODATION_SERVICE_DB_NAME"),
		DBPort:                             os.Getenv("ACCOMMODATION_SERVICE_DB_PORT"),
		ReservationServiceHost:             os.Getenv("RESERVATION_SERVICE_HOST"),
		ReservationServicePort:             os.Getenv("RESERVATION_SERVICE_PORT"),
		NatsHost:                           os.Getenv("NATS_HOST"),
		NatsPort:                           os.Getenv("NATS_PORT"),
		NatsUser:                           os.Getenv("NATS_USER"),
		NatsPass:                           os.Getenv("NATS_PASS"),
		DeleteUserCommandSubject:           os.Getenv("DELETE_USER_COMMAND_SUBJECT"),
		DeleteUserReplySubject:             os.Getenv("DELETE_USER_REPLY_SUBJECT"),
		AccommodationEventDbName:           os.Getenv("ACCOMMODATION_EVENT_DB_NAME"),
		AccommodationEventDbPort:           os.Getenv("ACCOMMODATION_EVENT_DB_PORT"),
		AccommodationEventInnerDbName:      os.Getenv("ACCOMMODATION_EVENT_INNER_DB_NAME"),
		AccommodationEventDbCollectionName: os.Getenv("ACCOMMODATION_EVENT_DB_COLLECTION_NAME"),
	}
}
