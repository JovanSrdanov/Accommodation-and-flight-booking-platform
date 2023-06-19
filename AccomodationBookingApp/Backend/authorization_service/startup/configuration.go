package startup

import "os"

type Configuration struct {
	Port                               string
	DBHost                             string
	DBPort                             string
	DBName                             string
	DBUser                             string
	DBPass                             string
	SecretKey                          string
	NatsHost                           string
	NatsPort                           string
	NatsUser                           string
	NatsPass                           string
	DeleteUserCommandSubject           string
	DeleteUserReplySubject             string
	AuthorizationEventDbName           string
	AuthorizationEventDbPort           string
	AuthorizationEventInnerDbName      string
	AuthorizationEventDbCollectionName string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                               os.Getenv("AUTHORIZATION_SERVICE_PORT"),
		DBHost:                             os.Getenv("AUTHORIZATION_DB_HOST"),
		DBPort:                             os.Getenv("AUTHORIZATION_DB_PORT"),
		DBName:                             os.Getenv("AUTHORIZATION_DB_NAME"),
		DBUser:                             os.Getenv("AUTHORIZATION_DB_USER"),
		DBPass:                             os.Getenv("AUTHORIZATION_DB_PASS"),
		SecretKey:                          os.Getenv("TOKEN_SYMMETRIC_KEY"),
		NatsHost:                           os.Getenv("NATS_HOST"),
		NatsPort:                           os.Getenv("NATS_PORT"),
		NatsUser:                           os.Getenv("NATS_USER"),
		NatsPass:                           os.Getenv("NATS_PASS"),
		DeleteUserCommandSubject:           os.Getenv("DELETE_USER_COMMAND_SUBJECT"),
		DeleteUserReplySubject:             os.Getenv("DELETE_USER_REPLY_SUBJECT"),
		AuthorizationEventDbName:           os.Getenv("AUTHORIZATION_EVENT_DB_NAME"),
		AuthorizationEventDbPort:           os.Getenv("AUTHORIZATION_EVENT_DB_PORT"),
		AuthorizationEventInnerDbName:      os.Getenv("AUTHORIZATION_EVENT_INNER_DB_NAME"),
		AuthorizationEventDbCollectionName: os.Getenv("AUTHORIZATION_EVENT_DB_COLLECTION_NAME"),
	}
}
