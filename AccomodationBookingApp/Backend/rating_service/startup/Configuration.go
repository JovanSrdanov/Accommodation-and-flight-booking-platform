package startup

import "os"

type Configuration struct {
	Port   string
	DbUser string
	DbPass string
	DbUri  string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:   os.Getenv("RATING_SERVICE_PORT"),
		DbUser: os.Getenv("NEO4J_USERNAME"),
		DbPass: os.Getenv("NEO4J_PASS"),
		DbUri:  os.Getenv("NEO4J_DB"),
	}
}
