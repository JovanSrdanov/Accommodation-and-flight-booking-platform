package startup

import "os"

type Configuration struct {
	Port            string
	DBName          string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	NatsHost        string
	NatsPort        string
	NatsUser        string
	NatsPass        string
	AuthServiceHost string
	AuthServicePort string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:            os.Getenv("NOTIFICATION_SERVICE_PORT"),
		DBName:          os.Getenv("NOTIFICATION_DB_NAME"),
		DBHost:          os.Getenv("NOTIFICATION_DB_HOST"),
		DBPort:          os.Getenv("NOTIFICATION_DB_PORT"),
		DBUser:          os.Getenv("NOTIFICATION_DB_USER"),
		DBPass:          os.Getenv("NOTIFICATION_DB_PASS"),
		NatsHost:        os.Getenv("NATS_HOST"),
		NatsPort:        os.Getenv("NATS_PORT"),
		NatsUser:        os.Getenv("NATS_USER"),
		NatsPass:        os.Getenv("NATS_PASS"),
		AuthServiceHost: os.Getenv("AUTHORIZATION_SERVICE_HOST"),
		AuthServicePort: os.Getenv("AUTHORIZATION_SERVICE_PORT"),
	}
}
