package startup

import "os"

type Configuration struct {
	Port   string
	DBName string
	DBPort string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:   os.Getenv("RESERVATION_SERVICE_PORT"),
		DBName: os.Getenv("RESERVATION_SERVICE_DB_NAME"),
		DBPort: os.Getenv("RESERVATION_SERVICE_DB_PORT"),
	}
}
