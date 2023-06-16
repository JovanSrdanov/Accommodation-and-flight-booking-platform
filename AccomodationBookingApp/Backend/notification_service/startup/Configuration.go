package startup

import "os"

type Configuration struct {
	Port                                  string
	DBName                                string
	DBHost                                string
	DBPort                                string
	DBUser                                string
	DBPass                                string
	NatsHost                              string
	NatsPort                              string
	NatsUser                              string
	NatsPass                              string
	AuthServiceHost                       string
	AuthServicePort                       string
	SendEventToNotificationServiceSubject string
	SendNotificationToAPIGatewaySubject   string
	NotificationEventDbName               string
	NotificationEventDbPort               string
	NotificationEventInnerDbName          string
	NotificationEventDbCollectionName     string
	DeleteUserCommandSubject              string
	DeleteUserReplySubject                string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                                  os.Getenv("NOTIFICATION_SERVICE_PORT"),
		DBName:                                os.Getenv("NOTIFICATION_DB_NAME"),
		DBHost:                                os.Getenv("NOTIFICATION_DB_HOST"),
		DBPort:                                os.Getenv("NOTIFICATION_DB_PORT"),
		DBUser:                                os.Getenv("NOTIFICATION_DB_USER"),
		DBPass:                                os.Getenv("NOTIFICATION_DB_PASS"),
		NatsHost:                              os.Getenv("NATS_HOST"),
		NatsPort:                              os.Getenv("NATS_PORT"),
		NatsUser:                              os.Getenv("NATS_USER"),
		NatsPass:                              os.Getenv("NATS_PASS"),
		AuthServiceHost:                       os.Getenv("AUTHORIZATION_SERVICE_HOST"),
		AuthServicePort:                       os.Getenv("AUTHORIZATION_SERVICE_PORT"),
		SendEventToNotificationServiceSubject: os.Getenv("SEND_EVENT_TO_NOTIFICATION_SERVICE_SUBJECT"),
		SendNotificationToAPIGatewaySubject:   os.Getenv("SEND_NOTIFICATION_TO_API_GATEWAY_SUBJECT"),
		NotificationEventDbName:               os.Getenv("NOTIFICATION_EVENT_DB_NAME"),
		NotificationEventDbPort:               os.Getenv("NOTIFICATION_EVENT_DB_PORT"),
		NotificationEventInnerDbName:          os.Getenv("NOTIFICATION_EVENT_INNER_DB_NAME"),
		NotificationEventDbCollectionName:     os.Getenv("NOTIFICATION_EVENT_DB_COLLECTION_NAME"),
		DeleteUserCommandSubject:              os.Getenv("DELETE_USER_COMMAND_SUBJECT"),
		DeleteUserReplySubject:                os.Getenv("DELETE_USER_REPLY_SUBJECT"),
	}
}
