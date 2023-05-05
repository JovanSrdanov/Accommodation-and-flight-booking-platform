package configuration

import "os"

type Configuration struct {
	Port      string
	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPass    string
	SecretKey string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:      os.Getenv("AUTHORIZATION_SERVICE_PORT"),
		DBHost:    os.Getenv("AUTHORIZATION_DB_HOST"),
		DBPort:    os.Getenv("AUTHORIZATION_DB_PORT"),
		DBName:    os.Getenv("AUTHORIZATION_DB_NAME"),
		DBUser:    os.Getenv("AUTHORIZATION_DB_USER"),
		DBPass:    os.Getenv("AUTHORIZATION_DB_PASS"),
		SecretKey: os.Getenv("TOKEN_SYMMETRIC_KEY"),
	}
}
