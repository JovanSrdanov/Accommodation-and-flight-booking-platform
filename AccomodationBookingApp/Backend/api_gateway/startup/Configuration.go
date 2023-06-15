package startup

import (
	"os"
)

type Configuration struct {
	Port                                string
	AuthorizationHost                   string
	AuthorizationPort                   string
	UserProfileHost                     string
	UserProfilePort                     string
	AccommodationHost                   string
	AccommodationPort                   string
	ReservationHost                     string
	ReservationPort                     string
	NotificationHost                    string
	NotificationPort                    string
	RatingHost                          string
	RatingPort                          string
	SendNotificationToAPIGatewaySubject string
	NatsHost                            string
	NatsPort                            string
	NatsUser                            string
	NatsPass                            string
}

func NewConfig() *Configuration {
	return &Configuration{
		Port:                                os.Getenv("GATEWAY_PORT"),
		AuthorizationHost:                   os.Getenv("AUTHORIZATION_SERVICE_HOST"),
		AuthorizationPort:                   os.Getenv("AUTHORIZATION_SERVICE_PORT"),
		UserProfileHost:                     os.Getenv("USER_PROFILE_SERVICE_HOST"),
		UserProfilePort:                     os.Getenv("USER_PROFILE_SERVICE_PORT"),
		AccommodationHost:                   os.Getenv("ACCOMMODATION_SERVICE_HOST"),
		AccommodationPort:                   os.Getenv("ACCOMMODATION_SERVICE_PORT"),
		ReservationHost:                     os.Getenv("RESERVATION_SERVICE_HOST"),
		ReservationPort:                     os.Getenv("RESERVATION_SERVICE_PORT"),
		NotificationHost:                    os.Getenv("NOTIFICATION_SERVICE_HOST"),
		NotificationPort:                    os.Getenv("NOTIFICATION_SERVICE_PORT"),
		RatingHost:                          os.Getenv("RATING_SERVICE_HOST"),
		RatingPort:                          os.Getenv("RATING_SERVICE_PORT"),
		SendNotificationToAPIGatewaySubject: os.Getenv("SEND_NOTIFICATION_TO_API_GATEWAY_SUBJECT"),
		NatsHost:                            os.Getenv("NATS_HOST"),
		NatsPort:                            os.Getenv("NATS_PORT"),
		NatsUser:                            os.Getenv("NATS_USER"),
		NatsPass:                            os.Getenv("NATS_PASS"),
	}
}
