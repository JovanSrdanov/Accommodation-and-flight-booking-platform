package startup

import "os"

type Configuration struct {
	Port                                  string
	DBName                                string
	DBPort                                string
	NatsHost                              string
	NatsPort                              string
	NatsUser                              string
	NatsPass                              string
	DeleteUserCommandSubject              string
	DeleteUserReplySubject                string
	SendEventToNotificationServiceSubject string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                                  os.Getenv("RESERVATION_SERVICE_PORT"),
		DBName:                                os.Getenv("RESERVATION_SERVICE_DB_NAME"),
		DBPort:                                os.Getenv("RESERVATION_SERVICE_DB_PORT"),
		NatsHost:                              os.Getenv("NATS_HOST"),
		NatsPort:                              os.Getenv("NATS_PORT"),
		NatsUser:                              os.Getenv("NATS_USER"),
		NatsPass:                              os.Getenv("NATS_PASS"),
		DeleteUserCommandSubject:              os.Getenv("DELETE_USER_COMMAND_SUBJECT"),
		DeleteUserReplySubject:                os.Getenv("DELETE_USER_REPLY_SUBJECT"),
		SendEventToNotificationServiceSubject: os.Getenv("SEND_EVENT_TO_NOTIFICATION_SERVICE_SUBJECT"),
	}
}
