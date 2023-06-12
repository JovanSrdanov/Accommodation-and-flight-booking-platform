package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"rating_service/domain/model"
)

type RatingRepositoryNeo4J struct {
	dbClient neo4j.Driver
}

func NewRatingRepositoryNeo4J(dbClient neo4j.Driver) (*RatingRepositoryNeo4J, error) {
	return &RatingRepositoryNeo4J{dbClient: dbClient}, nil
}

func (repo RatingRepositoryNeo4J) RateAccommodation(ratingg *model.Rating) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	guestID := "guest123"
	accommodationID := ratingg.AccommodationId.Hex()
	rating := ratingg.Rating
	date := ratingg.Date.String()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Check if Accommodation exists
		result, err := tx.Run(
			"MATCH (a:Accommodation {accommodationId: $accommodationID}) RETURN a",
			map[string]interface{}{
				"accommodationID": accommodationID,
			},
		)
		if err != nil {
			return nil, err
		}

		if !result.Next() {
			// Accommodation doesn't exist, create it
			_, err = tx.Run(
				"CREATE (a:Accommodation {accommodationId: $accommodationID})",
				map[string]interface{}{
					"accommodationID": accommodationID,
				},
			)
			if err != nil {
				return nil, err
			}
		}

		// Check if Guest exists
		result, err = tx.Run(
			"MATCH (g:Guest {guestId: $guestID}) RETURN g",
			map[string]interface{}{
				"guestID": guestID,
			},
		)
		if err != nil {
			return nil, err
		}

		if !result.Next() {
			// Guest doesn't exist, create it
			_, err = tx.Run(
				"CREATE (g:Guest {guestId: $guestID})",
				map[string]interface{}{
					"guestID": guestID,
				},
			)
			if err != nil {
				return nil, err
			}
		}

		// Check if relationship exists
		result, err = tx.Run(
			"MATCH (g:Guest {guestId: $guestID})-[r:RATED]->(a:Accommodation {accommodationId: $accommodationID}) WHERE r.date <> $date RETURN r",
			map[string]interface{}{
				"guestID":         guestID,
				"accommodationID": accommodationID,
				"date":            date,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			// Relationship exists, update it
			_, err = tx.Run(
				"MATCH (g:Guest {guestId: $guestID})-[r:RATED]->(a:Accommodation {accommodationId: $accommodationID}) SET r.rating = $rating, r.date = $date",
				map[string]interface{}{
					"guestID":         guestID,
					"accommodationID": accommodationID,
					"rating":          rating,
					"date":            date,
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			// Relationship doesn't exist, create it
			_, err = tx.Run(
				"MATCH (g:Guest {guestId: $guestID}), (a:Accommodation {accommodationId: $accommodationID}) CREATE (g)-[r:RATED {rating: $rating, date: $date}]->(a)",
				map[string]interface{}{
					"guestID":         guestID,
					"accommodationID": accommodationID,
					"rating":          rating,
					"date":            date,
				},
			)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (repo RatingRepositoryNeo4J) GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error) {
	return model.RatingResponse{
		AccommodationId: id.Hex(),
		Rating:          10,
	}, nil
}
