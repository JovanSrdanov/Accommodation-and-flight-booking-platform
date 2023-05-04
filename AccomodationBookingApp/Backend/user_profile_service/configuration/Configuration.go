package configuration

import "os"

type Configuration struct {
	Port   string
	DBName string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:   os.Getenv("USER_PROFILE_SERVICE_PORT"),
		DBName: os.Getenv("USER_PROFILE_DB_NAME"),
		DBHost: os.Getenv("USER_PROFILE_DB_HOST"),
		DBPort: os.Getenv("USER_PROFILE_DB_PORT"),
		DBUser: os.Getenv("USER_PROFILE_DB_USER"),
		DBPass: os.Getenv("USER_PROFILE_DB_PASS"),
	}
}
