package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"rating_service/domain/model"
	"sort"
	"time"
)

type RatingRepositoryNeo4J struct {
	dbClient neo4j.Driver
}

func NewRatingRepositoryNeo4J(dbClient neo4j.Driver) (*RatingRepositoryNeo4J, error) {
	return &RatingRepositoryNeo4J{dbClient: dbClient}, nil
}

func (repo RatingRepositoryNeo4J) RateAccommodation(guestId string, ratingDto *model.Rating) error {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	guestID := guestId
	accommodationID := ratingDto.AccommodationId.Hex()
	rating := ratingDto.Rating
	date := ratingDto.Date.Format("2006-01-02")

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
			"MATCH (g:Guest {guestId: $guestID})-[r:RATED]->(a:Accommodation {accommodationId: $accommodationID}) RETURN r",
			map[string]interface{}{
				"guestID":         guestID,
				"accommodationID": accommodationID,
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
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	accommodationID := id.Hex()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Retrieve ratings for the accommodation
		result, err := tx.Run(
			"MATCH (g:Guest)-[r:RATED]->(a:Accommodation {accommodationId: $accommodationID}) RETURN r.rating, g.guestId, r.date",
			map[string]interface{}{
				"accommodationID": accommodationID,
			},
		)
		if err != nil {
			return nil, err
		}

		ratingSum := int64(0)
		count := int64(0)
		singleRatings := make([]*model.SingleRating, 0)

		for result.Next() {
			record := result.Record()
			ratingValue, ok := record.Get("r.rating")
			guestId, ok := record.Get("g.guestId")
			dateStr, ok := record.Get("r.date")
			if ok {
				rating := ratingValue.(int64)
				ratingSum += rating
				count++
				date, _ := time.Parse("2006-01-02", dateStr.(string))
				singleRatings = append(singleRatings, &model.SingleRating{
					GuestId: guestId.(string),
					Rating:  int32(ratingValue.(int64)),
					Date:    date,
				})
			}
		}

		if count > 0 {
			ratingAverage := float64(ratingSum) / float64(count)

			res := model.RatingResponse{
				AccommodationId: accommodationID,
				AvgRating:       float32(ratingAverage),
				Ratings:         singleRatings,
			}

			return res, nil
		}

		return nil, nil
	})

	if err != nil {
		log.Fatal(err)
	}
	res := model.RatingResponse{}
	if result != nil {
		res = result.(model.RatingResponse)
	}

	return res, nil
}

func (repo RatingRepositoryNeo4J) GetRecommendedAccommodations(guestId string) ([]model.RecommendedAccommodation, error) {
	/*slice := make([]primitive.ObjectID, 0)
	slice = append(slice, primitive.NewObjectID())
	slice = append(slice, primitive.NewObjectID())
	return model.RecommendedAccommodation{AccommodationsIds: slice}, nil*/

	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	guestID := guestId

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Find similar guests
		result, err := tx.Run(
			"MATCH (g1:Guest {guestId: $guestID})-[r1:RATED]->(a:Accommodation)<-[r2:RATED]-(g2:Guest) "+
				"WHERE abs(r1.rating - r2.rating) <= 1 "+
				"WITH g1, g2, COUNT(DISTINCT a) AS commonRatings "+
				"WHERE g1 <> g2 AND commonRatings > 0 "+
				"RETURN g2.guestId",
			map[string]interface{}{
				"guestID": guestID,
			},
		)
		if err != nil {
			return nil, err
		}

		similarGuests := make([]string, 0)
		for result.Next() {
			record := result.Record()
			guestID, ok := record.Get("g2.guestId")
			if ok {
				similarGuests = append(similarGuests, guestID.(string))
			}
		}

		// Find accommodations rated by similar guests with rating >= 4
		result, err = tx.Run(
			"MATCH (:Guest {guestId: $guestID})-[:RATED]->(a1:Accommodation)<-[:RATED]-(g:Guest)-[r:RATED]->(a2:Accommodation) "+
				"WHERE a1 <> a2 AND g.guestId IN $similarGuests AND r.rating >= 4 "+
				"RETURN DISTINCT a2.accommodationId",
			map[string]interface{}{
				"guestID":       guestID,
				"similarGuests": similarGuests,
			},
		)
		if err != nil {
			return nil, err
		}

		recommendations := make([]string, 0)
		for result.Next() {
			record := result.Record()
			accommodationID, ok := record.Get("a2.accommodationId")
			if ok {
				recommendations = append(recommendations, accommodationID.(string))
			}
		}

		result, err = tx.Run(
			"MATCH (a:Accommodation) "+
				"WHERE a.accommodationId IN $accommodationIDs "+
				"WITH a, SIZE([(a)<-[r:RATED]-() WHERE r.rating <= 3 | r]) AS ratedCount "+
				"WHERE ratedCount < 3 "+
				"RETURN a.accommodationId",
			map[string]interface{}{
				"accommodationIDs": recommendations,
			},
		)
		if err != nil {
			return nil, err
		}

		accommodations := make([]string, 0)
		for result.Next() {
			record := result.Record()
			accommodation, ok := record.Get("a.accommodationId")
			if ok {
				accommodations = append(accommodations, accommodation.(string))
			}
		}

		//Calculate rating for each accommodation
		unsortedResponse := make([]model.RecommendedAccommodation, 0)
		for _, val := range accommodations {
			avgRat, _ := repo.CalculateRatingForAccommodationLocal(val)
			unsortedResponse = append(unsortedResponse, model.RecommendedAccommodation{
				AccommodationsId: val,
				AvgRating:        float32(avgRat),
			})
		}

		sort.Slice(unsortedResponse, func(i, j int) bool {
			return unsortedResponse[i].AvgRating > unsortedResponse[j].AvgRating
		})

		// Limit the result to the first ten elements if the slice is larger
		if len(unsortedResponse) > 5 {
			unsortedResponse = unsortedResponse[:5]
		}

		return unsortedResponse, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return result.([]model.RecommendedAccommodation), nil
}

func (repo RatingRepositoryNeo4J) DeleteRatingForAccommodation(accommodationId string, guestId string) (string, error) {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (:Guest {guestId: $guestID})-[r:RATED]->(:Accommodation {accommodationId: $accommodationID}) "+
				"DELETE r",
			map[string]interface{}{
				"guestID":         guestId,
				"accommodationID": accommodationId,
			},
		)
		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	if err != nil {
		log.Fatal(err)
	}

	return "Rating deleted", nil
}

func (repo RatingRepositoryNeo4J) RateHost(ratingDto *model.RateHostDto) error {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	guestID := ratingDto.GuestId
	hostID := ratingDto.HostId
	rating := ratingDto.Rating
	date := ratingDto.Date.Format("2006-01-02")

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Check if Accommodation exists
		result, err := tx.Run(
			"MATCH (a:Host {hostId: $hostID}) RETURN a",
			map[string]interface{}{
				"hostID": hostID,
			},
		)
		if err != nil {
			return nil, err
		}

		if !result.Next() {
			// Accommodation doesn't exist, create it
			_, err = tx.Run(
				"CREATE (a:Host {hostId: $hostID})",
				map[string]interface{}{
					"hostID": hostID,
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
			"MATCH (g:Guest {guestId: $guestID})-[r:RATED_HOST]->(a:Host {hostId: $hostID}) RETURN r",
			map[string]interface{}{
				"guestID": guestID,
				"hostID":  hostID,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			// Relationship exists, update it
			_, err = tx.Run(
				"MATCH (g:Guest {guestId: $guestID})-[r:RATED_HOST]->(a:Host {hostId: $hostID}) SET r.rating = $rating, r.date = $date",
				map[string]interface{}{
					"guestID": guestID,
					"hostID":  hostID,
					"rating":  rating,
					"date":    date,
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			// Relationship doesn't exist, create it
			_, err = tx.Run(
				"MATCH (g:Guest {guestId: $guestID}), (a:Host {hostId: $hostID}) CREATE (g)-[r:RATED_HOST {rating: $rating, date: $date}]->(a)",
				map[string]interface{}{
					"guestID": guestID,
					"hostID":  hostID,
					"rating":  rating,
					"date":    date,
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

func (repo RatingRepositoryNeo4J) GetRatingForHost(hostID string) (model.HostRatingResponse, error) {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Retrieve ratings for the accommodation
		result, err := tx.Run(
			"MATCH (g:Guest)-[r:RATED_HOST]->(a:Host {hostId: $hostID}) RETURN r.rating, g.guestId, r.date",
			map[string]interface{}{
				"hostID": hostID,
			},
		)
		if err != nil {
			return nil, err
		}

		ratingSum := int64(0)
		count := int64(0)
		singleRatings := make([]*model.SingleRating, 0)

		for result.Next() {
			record := result.Record()
			ratingValue, ok := record.Get("r.rating")
			guestId, ok := record.Get("g.guestId")
			dateStr, ok := record.Get("r.date")
			if ok {
				rating := ratingValue.(int64)
				ratingSum += rating
				count++
				date, _ := time.Parse("2006-01-02", dateStr.(string))
				singleRatings = append(singleRatings, &model.SingleRating{
					GuestId: guestId.(string),
					Rating:  int32(ratingValue.(int64)),
					Date:    date,
				})
			}
		}

		if count > 0 {
			ratingAverage := float64(ratingSum) / float64(count)

			res := model.HostRatingResponse{
				HostId:    hostID,
				AvgRating: float32(ratingAverage),
				Ratings:   singleRatings,
			}

			return res, nil
		}

		return nil, nil
	})

	if err != nil {
		log.Fatal(err)
	}
	res := model.HostRatingResponse{}
	if result != nil {
		res = result.(model.HostRatingResponse)
	}

	return res, nil
}

func (repo RatingRepositoryNeo4J) DeleteRatingForHost(hostId string, guestId string) (string, error) {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (:Guest {guestId: $guestID})-[r:RATED_HOST]->(:Host {hostId: $hostID}) "+
				"DELETE r",
			map[string]interface{}{
				"guestID": guestId,
				"hostID":  hostId,
			},
		)
		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	if err != nil {
		log.Fatal(err)
	}

	return "Rating for host deleted", nil
}

func (repo RatingRepositoryNeo4J) CalculateRatingForAccommodationLocal(accommodationID string) (float64, error) {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (:Accommodation {accommodationId: $accommodationID})<-[r:RATED]-(g:Guest) "+
				"RETURN toFloat(SUM(r.rating)) / count(r) AS avgRating",
			map[string]interface{}{
				"accommodationID": accommodationID,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			record := result.Record()
			avgRating, ok := record.Get("avgRating")
			if ok {
				return avgRating, nil
			}
		}

		return 0.0, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return result.(float64), nil
}

func (repo RatingRepositoryNeo4J) CalculateRatingForAccommodation(accommodationId string) (model.SimpleRatingResponse, error) {
	rating, err := repo.CalculateRatingForAccommodationLocal(accommodationId)
	if err != nil {
		return model.SimpleRatingResponse{}, err
	}
	return model.SimpleRatingResponse{
		AccommodationId: accommodationId,
		AvgRating:       float32(rating),
	}, nil
}

func (repo RatingRepositoryNeo4J) CalculateRatingForHost(hostId string) (model.SimpleHostRatingResponse, error) {
	session := repo.dbClient.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (:Host {hostId: $hostID})<-[r:RATED_HOST]-(g:Guest) "+
				"RETURN toFloat(SUM(r.rating)) / count(r) AS avgRating",
			map[string]interface{}{
				"hostID": hostId,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			record := result.Record()
			avgRating, ok := record.Get("avgRating")
			if ok {
				return avgRating, nil
			}
		}

		return 0.0, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return model.SimpleHostRatingResponse{
		HostId:    hostId,
		AvgRating: float32(result.(float64)),
	}, nil
}
