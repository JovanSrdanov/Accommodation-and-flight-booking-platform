package startup

import "os"

type Configuration struct {
	Port string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port: os.Getenv("ACCOMMODATION_SERVICE_PORT"),
	}
}
