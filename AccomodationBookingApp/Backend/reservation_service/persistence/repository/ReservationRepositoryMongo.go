package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"reservation_service/domain/model"
	"time"
)

type ReservationRepositoryMongo struct {
	dbClient *mongo.Client
}

func NewReservationRepositoryMongo(dbClient *mongo.Client) (*ReservationRepositoryMongo, error) {
	return &ReservationRepositoryMongo{dbClient: dbClient}, nil
}

/*
	CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error)
	GetAllMy() (model.Availabilities, error)
	UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error)
	CreateReservation(reservation *model.Reservation) (*model.Reservation, error)
	GetAllPendingReservations() (*model.Reservation, error)
	GetAllAcceptedReservations() (*model.Reservation, error)
	RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error)
*/

func (repo ReservationRepositoryMongo) CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error) {
	collection := repo.getCollectionAvailability()

	availability.PriceWithDate.ID = primitive.NewObjectID()

	filter := bson.D{{"accommodationId", availability.AccommodationId}}
	update := bson.M{"$push": bson.M{"availableDates": availability.PriceWithDate}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
	return primitive.ObjectID{}, nil
}

func (repo ReservationRepositoryMongo) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()

	filter := bson.D{{"accommodationId", reservation.AccommodationId}}

	// Find the availability document
	var availability model.Availability
	err := collection.FindOne(ctx, filter).Decode(&availability)
	if err != nil {
		log.Fatal(err)
		return &model.Reservation{}, err
	}

	//Da li postoji uopste range u koji pripada
	isRangeFound := false
	var foundRange *model.PriceWithDate

	for _, availableRange := range availability.AvailableDates {
		if reservation.DateRange.IsInside(availableRange.DateRange) {
			isRangeFound = true
			foundRange = availableRange
			break
		}
	}

	if !isRangeFound {
		return &model.Reservation{}, status.Errorf(codes.Aborted, "Not available date")
	}

	//Da li se preklapa sa nekom accpeted rezervacijom
	collectionReservations := repo.getCollectionReservation()

	filter2 := bson.D{{"accommodationId", reservation.AccommodationId}}
	cursor, err := collectionReservations.Find(ctx, filter2)
	if err != nil {
		log.Println(err)
		return &model.Reservation{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return &model.Reservation{}, err
	}

	for _, reservationValue := range reservations {
		if reservationValue.Status == "accepted" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			return &model.Reservation{}, status.Errorf(codes.Aborted, "Not available date, overlaps*")
		}
	}

	if availability.IsAutomaticReservation {
		reservation.Status = "accepted"
	} else {
		reservation.Status = "pending"
	}

	reservation.ID = primitive.NewObjectID()

	if foundRange.IsPricePerPerson {
		reservation.Price = foundRange.Price * reservation.NumberOfGuests
	} else {
		reservation.Price = foundRange.Price
	}

	_, err = collectionReservations.InsertOne(ctx, &reservation)
	if err != nil {
		log.Println(err)
		return &model.Reservation{}, err
	}

	return reservation, nil
}

func (repo ReservationRepositoryMongo) UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error) {
	//Nadji avail sa accId
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()

	filter := bson.D{{"accommodationId", priceWithDate.AccommodationId}}

	var availability model.Availability
	err := collection.FindOne(ctx, filter).Decode(&availability)
	if err != nil {
		log.Fatal(err)
		return &model.UpdatePriceAndDate{}, err
	}

	//U njemu nadji availableDate koji ima id kao prosledjeni
	var oldAvailableDate *model.PriceWithDate
	dateFound := false

	for _, availableDatesIterator := range availability.AvailableDates {
		if availableDatesIterator.ID == priceWithDate.PriceWithDate.ID {
			oldAvailableDate = availableDatesIterator
			dateFound = true
			break
		}
	}

	if !dateFound {
		return &model.UpdatePriceAndDate{}, status.Errorf(codes.Aborted, "This date doesn't exist in database wrong id")
	}

	//Sad proveri da li se dateRange novog availability-a poklapa sa nekim starim osim sam sa sobom
	for _, availableDatesIterator := range availability.AvailableDates {
		if availableDatesIterator.ID != priceWithDate.PriceWithDate.ID && availableDatesIterator.DateRange.Overlaps(priceWithDate.PriceWithDate.DateRange) {
			return &model.UpdatePriceAndDate{}, status.Errorf(codes.Aborted, "Provided date overlaps with existing available date")
		}
	}

	//Proveri da li se taj availDate poklapa sa nekom od rezervacija koja je accepted ili pending
	collectionReservations := repo.getCollectionReservation()

	filter2 := bson.D{{"accommodationId", priceWithDate.AccommodationId}}
	cursor, err := collectionReservations.Find(ctx, filter2)
	if err != nil {
		log.Println(err)
		return &model.UpdatePriceAndDate{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return &model.UpdatePriceAndDate{}, err
	}

	for _, reservationValue := range reservations {
		if (reservationValue.Status == "accepted" || reservationValue.Status == "pending") && oldAvailableDate.DateRange.Overlaps(reservationValue.DateRange) {
			return &model.UpdatePriceAndDate{}, status.Errorf(codes.Aborted, "Can not modify this availability because it overlaps")
		}
	}

	//Ako se ne poklapa izmeni datum i cenu
	filter = bson.D{{"availableDates._id", priceWithDate.PriceWithDate.ID}}

	update := bson.M{
		"$set": bson.M{
			"availableDates.$.price":            priceWithDate.PriceWithDate.Price,
			"availableDates.$.isPricePerPerson": priceWithDate.PriceWithDate.IsPricePerPerson,
			"availableDates.$.dateRange": bson.M{
				"from": priceWithDate.PriceWithDate.DateRange.From,
				"to":   priceWithDate.PriceWithDate.DateRange.To,
			},
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return &model.UpdatePriceAndDate{}, err
	}

	return &model.UpdatePriceAndDate{}, nil
}

func (repo ReservationRepositoryMongo) GetAllMy(hostId string) (model.Availabilities, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()

	filter := bson.D{{"hostId", hostId}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return model.Availabilities{}, err
	}

	var availabilities model.Availabilities
	err = cursor.All(ctx, &availabilities)
	if err != nil {
		log.Println(err)
		return model.Availabilities{}, err
	}
	return availabilities, nil
}

func (repo ReservationRepositoryMongo) GetAllAcceptedReservations(hostId string) (model.Reservations, error) {
	//Dobavi sve dostpunosti i iz njih izvuci sve accommodationId gde je hostId prosledjenji
	availabilities, err := repo.GetAllMy(hostId)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	accommodationIds := make([]primitive.ObjectID, 0)

	for _, availability := range availabilities {
		accommodationIds = append(accommodationIds, availability.AccommodationId)
	}

	//Sad dobavi sve rezervacije koje imaju ovaj accId i koje su accepted
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.M{
		"accommodationId": bson.M{"$in": accommodationIds},
		"status":          "accepted"}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	return reservations, nil
}

func (repo ReservationRepositoryMongo) GetAllPendingReservations(hostId string) (model.Reservations, error) {
	//Dobavi sve dostpunosti i iz njih izvuci sve accommodationId gde je hostId prosledjenji
	availabilities, err := repo.GetAllMy(hostId)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	accommodationIds := make([]primitive.ObjectID, 0)

	for _, availability := range availabilities {
		accommodationIds = append(accommodationIds, availability.AccommodationId)
	}

	//Sad dobavi sve rezervacije koje imaju ovaj accId i koje su accepted
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.M{
		"accommodationId": bson.M{"$in": accommodationIds},
		"status":          "pending"}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	return reservations, nil
}

func (repo ReservationRepositoryMongo) GetAllReservationsForGuest(guestId string) (model.Reservations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.D{{"guestId", guestId}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return model.Reservations{}, err
	}

	return reservations, nil
}

func (repo ReservationRepositoryMongo) CreateAvailabilityBase(base *model.Availability) (primitive.ObjectID, error) {
	base.AvailableDates = make([]*model.PriceWithDate, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()
	base.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(ctx, &base)
	if err != nil {
		log.Println(err)
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (repo ReservationRepositoryMongo) CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.D{{"_id", id}}

	// Find the reservation u want to change
	var reservation model.Reservation
	err := collection.FindOne(ctx, filter).Decode(&reservation)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}

	if reservation.Status == "accepted" || reservation.Status == "rejected" {
		return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "Can not cancel accepted or rejected reservation")
	}

	update := bson.M{"$set": bson.M{"status": "canceled"}}

	// Update the document matching the filter
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "Not updated")
	}

	return id, nil
}

func (repo ReservationRepositoryMongo) AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.D{{"_id", id}}

	// Find the reservation u want to change
	var reservation model.Reservation
	err := collection.FindOne(ctx, filter).Decode(&reservation)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}

	if reservation.Status != "pending" {
		return primitive.ObjectID{}, status.Errorf(codes.Aborted, "Can not accept reservation that is not pending")
	}

	//Find all reservation from same accommodation
	collectionReservations := repo.getCollectionReservation()

	filter = bson.D{{"accommodationId", reservation.AccommodationId}}
	cursor, err := collectionReservations.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return primitive.ObjectID{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		log.Println(err)
		return primitive.ObjectID{}, err
	}

	pendingIds := make([]primitive.ObjectID, 0)

	for _, reservationValue := range reservations {
		if reservationValue.ID != reservation.ID && reservationValue.Status == "accepted" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			return primitive.ObjectID{}, status.Errorf(codes.Aborted, "Can not accept this reservation, overlaps*")
		}
		if reservationValue.ID != reservation.ID && reservationValue.Status == "pending" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			pendingIds = append(pendingIds, reservationValue.ID)
		}
	}

	pendingIdString := pendingIds[0].String()
	log.Println(pendingIdString)

	//Ako se ne overlap onda accept
	filter = bson.D{{"_id", reservation.ID}}

	// Define an update to set the "status" field to "rejected"
	update := bson.M{"$set": bson.M{"status": "accepted"}}

	// Update the document matching the filter
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		return primitive.ObjectID{}, status.Errorf(codes.Canceled, "Not updated")
	}

	filter = bson.D{{"_id", bson.M{"$in": pendingIds}}}

	// Define an update to set the "status" field to "rejected"
	update = bson.M{"$set": bson.M{"status": "rejected"}}

	// Update the document matching the filter
	result, err = collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return reservation.ID, nil
}

func (repo ReservationRepositoryMongo) RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	//Ako se ne overlap onda accept
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.D{{"_id", id}}

	// Define an update to set the "status" field to "rejected"
	update := bson.M{"$set": bson.M{"status": "rejected"}}

	// Update the document matching the filter
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "Not updated")
	}

	return id, nil
}

func (repo ReservationRepositoryMongo) getCollectionAvailability() *mongo.Collection {
	db := repo.dbClient.Database("reservationDb")
	collection := db.Collection("availabilities")
	return collection
}

func (repo ReservationRepositoryMongo) getCollectionReservation() *mongo.Collection {
	db := repo.dbClient.Database("reservationDb")
	collection := db.Collection("reservations")
	return collection
}
