package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
)

func GetClient(dbUri string, dbUser string, dbPass string) (neo4j.Driver, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()

	auth := neo4j.BasicAuth(dbUser, dbPass, "")

	driver, err := neo4j.NewDriver(dbUri, auth)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	err = driver.VerifyConnectivity()
	if err != nil {
		log.Panic("Gascina " + err.Error())
		return nil, err
	}

	log.Printf(`Neo4J server address: %s`, driver.Target().Host)

	return driver, nil
}
