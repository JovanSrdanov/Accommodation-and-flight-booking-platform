package startup

import (
	"os"
)

type Configuration struct {
	Port              string
	AuthorizationHost string
	AuthorizationPort string
	UserProfileHost   string
	UserProfilePort   string
}

func NewConfig() *Configuration {
	return &Configuration{
		Port:              os.Getenv("GATEWAY_PORT"),
		AuthorizationHost: os.Getenv("AUTHORIZATION_SERVICE_HOST"),
		AuthorizationPort: os.Getenv("AUTHORIZATION_SERVICE_PORT"),
		UserProfileHost:   os.Getenv("USER_PROFILE_SERVICE_HOST"),
		UserProfilePort:   os.Getenv("USER_PROFILE_SERVICE_PORT"),
	}
}
