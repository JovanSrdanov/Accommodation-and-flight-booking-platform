package startup

import "os"

type Configuration struct {
	Port                                  string
	DbUser                                string
	DbPass                                string
	DbUri                                 string
	NatsHost                              string
	NatsPort                              string
	NatsUser                              string
	NatsPass                              string
	SendEventToNotificationServiceSubject string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                                  os.Getenv("RATING_SERVICE_PORT"),
		DbUser:                                os.Getenv("NEO4J_USERNAME"),
		DbPass:                                os.Getenv("NEO4J_PASS"),
		DbUri:                                 os.Getenv("NEO4J_DB"),
		NatsHost:                              os.Getenv("NATS_HOST"),
		NatsPort:                              os.Getenv("NATS_PORT"),
		NatsUser:                              os.Getenv("NATS_USER"),
		NatsPass:                              os.Getenv("NATS_PASS"),
		SendEventToNotificationServiceSubject: os.Getenv("SEND_EVENT_TO_NOTIFICATION_SERVICE_SUBJECT"),
	}
}
