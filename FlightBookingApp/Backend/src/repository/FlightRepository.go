package repository

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	utils "FlightBookingApp/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type flightRepository struct {
	base Repository
}

type FlightRepository interface {
	Create(flight *model.Flight) (primitive.ObjectID, error)
	GetAll() (model.Flights, error)
	GetById(id primitive.ObjectID) (model.Flight, error)
	Cancel(id primitive.ObjectID) error
	Search(flightSearchParameters *dto.FlightSearchParameters, pageInfo *utils.PageInfo) (*utils.Page, error)
}

// NoSQL: Constructor which reads db configuration from environment
func NewFlightRepository(client *mongo.Client, logger *log.Logger) *flightRepository {
	base := NewRepository(client, logger)
	return &flightRepository{base: base}
}

func (repo *flightRepository) Create(flight *model.Flight) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result, err := collection.InsertOne(ctx, &flight)
	if err != nil {
		repo.base.logger.Println(err)
		return primitive.ObjectID{}, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	repo.base.logger.Printf("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo *flightRepository) GetAll() (model.Flights, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	var flights model.Flights
	err = cursor.All(ctx, &flights)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return flights, nil
}

func (repo *flightRepository) GetById(id primitive.ObjectID) (model.Flight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.Flight{}, result.Err()
	}

	var flight model.Flight
	result.Decode(&flight)

	return flight, nil
}
func (repo *flightRepository) Cancel(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"canceled", true}}}}

	//Before calling Cancel, object is found so we know it will be found now and updated
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (repo *flightRepository) Search(flightSearchParameters *dto.FlightSearchParameters, pageInfo *utils.PageInfo) (*utils.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	filterOptions := filterSetup(flightSearchParameters)
	sortOptions, err := sortSetup(pageInfo)

	if err != nil {
		return nil, err
	}

	cursor, _ := collection.Find(ctx, filterOptions, sortOptions)

	var flights model.Flights
	err = cursor.All(ctx, &flights)
	if err != nil {
		return nil, err
	}

	flightSearchResults := createPageData(flightSearchParameters.DesiredNumberOfSeats, flights)
	count, err := collection.CountDocuments(ctx, filterOptions)
	if err != nil {
		return nil, err
	}
	page := utils.Page{flightSearchResults, int(count)}

	return &page, nil
}

func createPageData(DesiredNumberOfSeats int, flights model.Flights) []*dto.FlightSearchResult {
	var flightSearchResults []*dto.FlightSearchResult
	for _, element := range flights {
		flightSearchResults = append(flightSearchResults, dto.NewFlightSearchResult(element, DesiredNumberOfSeats))
	}
	return flightSearchResults
}

func sortSetup(pageInfo *utils.PageInfo) (*options.FindOptions, error) {
	if pageInfo.SortType != "time" && pageInfo.SortType != "price" {
		return nil, &errors.InvalidSortTypeError{}
	}
	sortDirection := 1
	if pageInfo.SortDirection == "dsc" {
		sortDirection = -1
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{pageInfo.SortType, sortDirection}})
	findOptions.SetSkip(int64(((pageInfo.PageNumber) - 1) * pageInfo.ResultsPerPage))
	findOptions.SetLimit(int64(pageInfo.ResultsPerPage))
	return findOptions, nil
}

func filterSetup(flightSearchParameters *dto.FlightSearchParameters) bson.M {
	filter := bson.M{
		"vacantSeats":                 bson.M{"$gte": flightSearchParameters.DesiredNumberOfSeats},
		"startPoint.address.city":     bson.M{"$regex": primitive.Regex{Pattern: flightSearchParameters.StartPointCity, Options: "i"}},
		"startPoint.address.country":  bson.M{"$regex": primitive.Regex{Pattern: flightSearchParameters.StartPointCountry, Options: "i"}},
		"destination.address.city":    bson.M{"$regex": primitive.Regex{Pattern: flightSearchParameters.DestinationCity, Options: "i"}},
		"destination.address.country": bson.M{"$regex": primitive.Regex{Pattern: flightSearchParameters.DestinationCountry, Options: "i"}},
		"time":                        bson.M{"$gte": flightSearchParameters.Date, "$lte": flightSearchParameters.Date.AddDate(0, 0, 1)},
	}
	return filter
}

func (repo *flightRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("flights")
	return collection
}
