package startup

import "os"

type Configuration struct {
	Port                   string
	DBName                 string
	DBPort                 string
	ReservationServiceHost string
	ReservationServicePort string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:                   os.Getenv("ACCOMMODATION_SERVICE_PORT"),
		DBName:                 os.Getenv("ACCOMMODATION_SERVICE_DB_NAME"),
		DBPort:                 os.Getenv("ACCOMMODATION_SERVICE_DB_PORT"),
		ReservationServiceHost: os.Getenv("RESERVATION_SERVICE_HOST"),
		ReservationServicePort: os.Getenv("RESERVATION_SERVICE_PORT"),
	}
}
