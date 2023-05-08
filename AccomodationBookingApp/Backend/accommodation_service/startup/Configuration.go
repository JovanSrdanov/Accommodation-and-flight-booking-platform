package startup

import "os"

type Configuration struct {
	Port   string
	DbUser string
	DbPass string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:   os.Getenv("ACCOMMODATION_SERVICE_PORT"),
		DbUser: os.Getenv("ACCOMMODATION_SERVICE_DB_USER"),
		DbPass: os.Getenv("ACCOMMODATION_SERVICE_DB_PASS"),
	}
}
