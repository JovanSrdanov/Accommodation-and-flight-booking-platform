package startup

import "os"

type Configuration struct {
	Port                     string
	DBName                   string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPass                   string
	NatsHost                 string
	NatsPort                 string
	NatsUser                 string
	NatsPass                 string
	DeleteUserCommandSubject string
	DeleteUserReplySubject   string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                     os.Getenv("USER_PROFILE_SERVICE_PORT"),
		DBName:                   os.Getenv("USER_PROFILE_DB_NAME"),
		DBHost:                   os.Getenv("USER_PROFILE_DB_HOST"),
		DBPort:                   os.Getenv("USER_PROFILE_DB_PORT"),
		DBUser:                   os.Getenv("USER_PROFILE_DB_USER"),
		DBPass:                   os.Getenv("USER_PROFILE_DB_PASS"),
		NatsHost:                 os.Getenv("NATS_HOST"),
		NatsPort:                 os.Getenv("NATS_PORT"),
		NatsUser:                 os.Getenv("NATS_USER"),
		NatsPass:                 os.Getenv("NATS_PASS"),
		DeleteUserCommandSubject: os.Getenv("DELETE_USER_COMMAND_SUBJECT"),
		DeleteUserReplySubject:   os.Getenv("DELETE_USER_REPLY_SUBJECT"),
	}
}
