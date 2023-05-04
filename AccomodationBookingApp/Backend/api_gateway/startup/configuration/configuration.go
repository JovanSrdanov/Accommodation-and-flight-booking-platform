package configuration

import (
	"os"
)

type Configuration struct {
	Port              string
	AuthorizationHost string
	AuthorizationPort string
}

func NewConfig() *Configuration {
	return &Configuration{
		Port:              os.Getenv("GATEWAY_PORT"),
		AuthorizationHost: os.Getenv("AUTHORIZATION_SERVICE_HOST"),
		AuthorizationPort: os.Getenv("AUTHORIZATION_SERVICE_PORT"),
	}
}
