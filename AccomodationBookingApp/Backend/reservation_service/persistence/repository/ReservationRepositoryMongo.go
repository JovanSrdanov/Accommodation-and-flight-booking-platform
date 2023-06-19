package repository

import (
	"common/NotificationMessaging"
	"common/saga/messaging"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reservation_service/domain/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReservationRepositoryMongo struct {
	dbClient  *mongo.Client
	publisher messaging.Publisher
}

func NewReservationRepositoryMongo(dbClient *mongo.Client, publisher messaging.Publisher) (*ReservationRepositoryMongo, error) {
	return &ReservationRepositoryMongo{
		dbClient:  dbClient,
		publisher: publisher,
	}, nil
}

func (repo ReservationRepositoryMongo) CreateAvailability(newAvailability *model.AvailabilityRequest) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()
	filter := bson.D{{"accommodationId", newAvailability.AccommodationId}}

	// Find the newAvailability document
	var availability model.Availability
	err := collection.FindOne(ctx, filter).Decode(&availability)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}

	for _, availableDate := range availability.AvailableDates {
		if availableDate.DateRange.Overlaps(newAvailability.PriceWithDate.DateRange) {
			return primitive.ObjectID{}, status.Errorf(codes.Aborted, "Can not define this availability, overlaps with existing one*")
		}
	}

	newAvailability.PriceWithDate.ID = primitive.NewObjectID()

	update := bson.M{"$push": bson.M{"availableDates": newAvailability.PriceWithDate}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return primitive.ObjectID{}, nil
}

func (repo ReservationRepositoryMongo) SearchAccommodation(accommodationIds []*primitive.ObjectID, dateRange model.DateRange, numberOfGuests int32) ([]*model.SearchResponseDto, error) {
	validAccommodationIds := make([]*model.SearchResponseDto, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()

	for _, accId := range accommodationIds {

		filter := bson.D{{"accommodationId", accId}}

		// Find the availability document
		var availability model.Availability
		err := collection.FindOne(ctx, filter).Decode(&availability)
		if err != nil {
			log.Fatal(err)
			return []*model.SearchResponseDto{}, err
		}

		//Da li moze da se spoje neke
		sortedAvailableDates := availability.AvailableDates
		bubbleSort(sortedAvailableDates)

		startFound := false
		endFound := false
		foundDates := make([]*model.PriceWithDate, 0)

		for _, date := range sortedAvailableDates {
			if dateRange.IsInside(date.DateRange) {
				startFound = true
				endFound = true
				foundDates = append(foundDates, date)
				break
			}

			if date.DateRange.IsStartFor(dateRange) {
				startFound = true
				foundDates = append(foundDates, date)
			} else if startFound && date.DateRange.Extends(foundDates[len(foundDates)-1].DateRange) {
				foundDates = append(foundDates, date)
				if date.DateRange.IsEndFor(dateRange) {
					endFound = true
					break
				}
			}
		}

		if !startFound || !endFound {
			//return []*model.SearchResponseDto{}, status.Errorf(codes.Aborted, "Not available date")
			continue
		}

		//Da li se preklapa sa nekom accpeted rezervacijom
		collectionReservations := repo.getCollectionReservation()

		cursor, err := collectionReservations.Find(ctx, filter)
		if err != nil {

			return []*model.SearchResponseDto{}, err
		}

		var reservations []*model.Reservation
		err = cursor.All(ctx, &reservations)
		if err != nil {

			return []*model.SearchResponseDto{}, err
		}

		gas := false

		for _, reservationValue := range reservations {

			if reservationValue.Status == "accepted" && reservationValue.DateRange.Overlaps(dateRange) {
				//return []*model.SearchResponseDto{}, status.Errorf(codes.Aborted, "Not available date, overlaps*")
				gas = true
				break
			}
		}

		if gas {
			continue
		}

		price := calculatePrice(foundDates, &model.Reservation{
			DateRange:      dateRange,
			NumberOfGuests: numberOfGuests,
		})

		validAccommodationIds = append(validAccommodationIds, &model.SearchResponseDto{
			AccommodationId: accId,
			Price:           price,
		})
	}

	return validAccommodationIds, nil
}

func (repo ReservationRepositoryMongo) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	accountID, err := uuid.Parse(availability.HostId)
	if err != nil {
		log.Fatal(err)
		return &model.Reservation{}, err
	}

	//Da li moze da se spoje neke
	sortedAvailableDates := availability.AvailableDates
	bubbleSort(sortedAvailableDates)

	startFound := false
	endFound := false
	foundDates := make([]*model.PriceWithDate, 0)

	for _, date := range sortedAvailableDates {
		if reservation.DateRange.IsInside(date.DateRange) {
			startFound = true
			endFound = true
			foundDates = append(foundDates, date)
			break
		}

		if date.DateRange.IsStartFor(reservation.DateRange) {
			startFound = true
			foundDates = append(foundDates, date)
		} else if startFound && date.DateRange.Extends(foundDates[len(foundDates)-1].DateRange) {
			foundDates = append(foundDates, date)
			if date.DateRange.IsEndFor(reservation.DateRange) {
				endFound = true
				break
			}
		}
	}

	if !startFound || !endFound {
		return &model.Reservation{}, status.Errorf(codes.Aborted, "Not available date")
	}

	//Da li se preklapa sa nekom accpeted rezervacijom
	collectionReservations := repo.getCollectionReservation()

	cursor, err := collectionReservations.Find(ctx, filter)
	if err != nil {
		return &model.Reservation{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return &model.Reservation{}, err
	}

	for _, reservationValue := range reservations {
		if reservationValue.Status == "accepted" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			return &model.Reservation{}, status.Errorf(codes.Aborted, "Not available date, overlaps*")
		}
	}

	oldProminenet, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return nil, err
	}
	//Kreiranje rezervacije

	if availability.IsAutomaticReservation {
		reservation.Status = "accepted"
		message := NotificationMessaging.NotificationMessage{
			MessageType:            "RequestMade",
			MessageForNotification: "A reservation has been made for your accommodation",
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
	} else {
		reservation.Status = "pending"
		message := NotificationMessaging.NotificationMessage{
			MessageType:            "RequestMade",
			MessageForNotification: "A reservation request has been made for your accommodation",
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
	}

	reservation.ID = primitive.NewObjectID()

	//!!! Kalkulisanje cene
	reservation.Price = calculatePrice(foundDates, reservation)

	_, err = collectionReservations.InsertOne(ctx, &reservation)
	if err != nil {
		return &model.Reservation{}, err
	}
	if availability.IsAutomaticReservation {
		accountID, err := uuid.Parse(reservation.GuestId)
		if err != nil {
			log.Fatal(err)
			return &model.Reservation{}, err
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "HostResponded",
			MessageForNotification: "The host has accepted your reservation request automatically",
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
	}
	newProminent, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return nil, err
	}
	if oldProminenet != newProminent {
		accountID, err := uuid.Parse(availability.HostId)
		if err != nil {
			log.Fatal(err)
			return &model.Reservation{}, err
		}
		MFN := ""
		if newProminent {
			MFN = "You are now a prominent host"
		} else {
			MFN = "You are no longer prominent host"
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "ProminentHost",
			MessageForNotification: MFN,
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
	}

	return reservation, nil
}

func calculatePrice(dates []*model.PriceWithDate, reservation *model.Reservation) int32 {
	var price int32 = 0
	for _, date := range dates {
		commonDays := date.DateRange.DaysInCommon(reservation.DateRange)
		if date.IsPricePerPerson {
			price += commonDays * date.Price * reservation.NumberOfGuests
		} else {
			price += commonDays * date.Price
		}
	}

	return price
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
		return &model.UpdatePriceAndDate{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
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
		return model.Availabilities{}, err
	}

	var availabilities model.Availabilities
	err = cursor.All(ctx, &availabilities)
	if err != nil {
		return model.Availabilities{}, err
	}
	return availabilities, nil
}

func (repo ReservationRepositoryMongo) GetAllAcceptedReservations(hostId string) (model.Reservations, error) {
	//Dobavi sve dostpunosti i iz njih izvuci sve accommodationId gde je hostId prosledjenji
	availabilities, err := repo.GetAllMy(hostId)
	if err != nil {
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
		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return model.Reservations{}, err
	}

	return reservations, nil
}

func (repo ReservationRepositoryMongo) GetAllPendingReservations(hostId string) (model.Reservations, []int32, error) {
	//Dobavi sve dostpunosti i iz njih izvuci sve accommodationId gde je hostId prosledjenji
	availabilities, err := repo.GetAllMy(hostId)
	if err != nil {
		return model.Reservations{}, []int32{}, err
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
		return model.Reservations{}, []int32{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return model.Reservations{}, []int32{}, err
	}

	cancelSlice := make([]int32, 0)

	for _, resr := range reservations {

		filter2 := bson.M{
			"guestId": resr.GuestId,
			"status":  "canceled"}

		count, err := collection.CountDocuments(ctx, filter2)
		if err != nil {
			return model.Reservations{}, []int32{}, err
		}

		cancelSlice = append(cancelSlice, int32(count))
	}

	return reservations, cancelSlice, nil
}

func (repo ReservationRepositoryMongo) GetAllReservationsForGuest(guestId string) (model.Reservations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.D{{"guestId", guestId}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {
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

	if reservation.Status == "rejected" {
		return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "Can not cancel rejected reservation")
	}

	update := bson.M{"$set": bson.M{"status": "canceled"}}

	filter = bson.D{{"accommodationId", reservation.AccommodationId}}

	coll := repo.getCollectionAvailability()
	// Find the availability document
	var availability model.Availability
	err = coll.FindOne(ctx, filter).Decode(&availability)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}
	oldProminenet, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	filter = bson.D{{"_id", id}}
	// Update the document matching the filter
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "Not updated")
	}
	newProminent, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if oldProminenet != newProminent {
		accountID, err := uuid.Parse(availability.HostId)
		if err != nil {
			log.Fatal(err)
			return primitive.ObjectID{}, err
		}
		MFN := ""
		if newProminent {
			MFN = "You are now a prominent host"
		} else {
			MFN = "You are no longer prominent host"
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "ProminentHost",
			MessageForNotification: MFN,
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
	}

	accountID, err := uuid.Parse(availability.HostId)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}
	message := NotificationMessaging.NotificationMessage{
		MessageType:            "ReservationCanceled",
		MessageForNotification: "Guest canceled a reservation",
		AccountID:              accountID,
	}
	repo.publisher.Publish(message)

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
		return primitive.ObjectID{}, err
	}

	var reservations []*model.Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	pendingIds := make([]primitive.ObjectID, 0)

	filter = bson.D{{"accommodationId", reservation.AccommodationId}}

	coll := repo.getCollectionAvailability()
	// Find the availability document
	var availability model.Availability
	err = coll.FindOne(ctx, filter).Decode(&availability)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}
	oldProminenet, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	for _, reservationValue := range reservations {
		if reservationValue.ID != reservation.ID && reservationValue.Status == "accepted" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			return primitive.ObjectID{}, status.Errorf(codes.Aborted, "Can not accept this reservation, overlaps*")
		}
		if reservationValue.ID != reservation.ID && reservationValue.Status == "pending" && reservationValue.DateRange.Overlaps(reservation.DateRange) {
			pendingIds = append(pendingIds, reservationValue.ID)
			accountID, err := uuid.Parse(reservationValue.GuestId)
			if err != nil {
				log.Fatal(err)
				return primitive.ObjectID{}, err
			}

			message := NotificationMessaging.NotificationMessage{
				MessageType:            "HostResponded",
				MessageForNotification: "Your reservation request has been rejected",
				AccountID:              accountID,
			}
			repo.publisher.Publish(message)

		}
	}

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

	/////////////////////

	accountID, err := uuid.Parse(reservation.GuestId)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}

	message := NotificationMessaging.NotificationMessage{
		MessageType:            "HostResponded",
		MessageForNotification: "The host has accepted your reservation request",
		AccountID:              accountID,
	}
	repo.publisher.Publish(message)

	newProminent, err := prominentHostHttp(availability.HostId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if oldProminenet != newProminent {
		accountID, err := uuid.Parse(availability.HostId)
		if err != nil {
			log.Fatal(err)
			return primitive.ObjectID{}, err
		}
		MFN := ""
		if newProminent {
			MFN = "You are now a prominent host"
		} else {
			MFN = "You are no longer prominent host"
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "ProminentHost",
			MessageForNotification: MFN,
			AccountID:              accountID,
		}
		repo.publisher.Publish(message)
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

	var reservation model.Reservation
	err = collection.FindOne(ctx, filter).Decode(&reservation)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}
	accountID, err := uuid.Parse(reservation.GuestId)
	if err != nil {
		log.Fatal(err)
		return primitive.ObjectID{}, err
	}

	message := NotificationMessaging.NotificationMessage{
		MessageType:            "HostResponded",
		MessageForNotification: "Your reservation request has been rejected",
		AccountID:              accountID,
	}
	repo.publisher.Publish(message)

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

func bubbleSort(nums []*model.PriceWithDate) {
	n := len(nums)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if nums[j].DateRange.From.After(nums[j+1].DateRange.From) {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
}

func (repo ReservationRepositoryMongo) GuestHasActiveReservations(guestID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"guestId":        guestID.String(),
		"dateRange.from": bson.M{"$gte": time.Now().UTC()},
		"status":         "accepted",
	}

	reservations := repo.getCollectionReservation()

	count, err := reservations.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (repo ReservationRepositoryMongo) DeleteAvailabilitiesAndReservationsByAccommodationId(accommodationId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"accommodationId": accommodationId,
	}

	reservations := repo.getCollectionReservation()
	availabilities := repo.getCollectionAvailability()

	_, err := reservations.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	_, err = availabilities.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo ReservationRepositoryMongo) GetAllReservationsForHost(hostId string) (model.Reservations, error) {
	//Dobavi sve dostpunosti i iz njih izvuci sve accommodationId gde je hostId prosledjenji
	availabilities, err := repo.GetAllMy(hostId)
	if err != nil {

		return model.Reservations{}, err
	}

	accommodationIds := make([]primitive.ObjectID, 0)

	for _, availability := range availabilities {
		accommodationIds = append(accommodationIds, availability.AccommodationId)
	}

	//Sad dobavi sve rezervacije koje imaju ovaj accId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.M{
		"accommodationId": bson.M{"$in": accommodationIds}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {

		return model.Reservations{}, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {

		return model.Reservations{}, err
	}

	return reservations, nil
}

func (repo ReservationRepositoryMongo) GetAllRatableAccommodationsForGuest(guestId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionReservation()

	filter := bson.M{
		"guestId": guestId,
		"status":  "accepted",
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {

		return nil, err
	}

	var reservations model.Reservations
	err = cursor.All(ctx, &reservations)
	if err != nil {

		return nil, err
	}

	resp := make([]string, 0)
	seen := make(map[string]bool)
	for _, val := range reservations {
		vreme := time.Now().In(time.UTC).Truncate(24 * time.Hour)
		if val.DateRange.To.Before(vreme) {
			accommodationID := val.AccommodationId.Hex()
			if !seen[accommodationID] {
				resp = append(resp, accommodationID)
				seen[accommodationID] = true
			}
		}
	}

	return resp, nil
}

func (repo ReservationRepositoryMongo) GetAllRatableHostsForGuest(guestId string) ([]string, error) {
	accommodationIds, err := repo.GetAllRatableAccommodationsForGuest(guestId)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionAvailability()

	accommodationOids := make([]primitive.ObjectID, 0)
	for _, val := range accommodationIds {
		oid, _ := primitive.ObjectIDFromHex(val)
		accommodationOids = append(accommodationOids, oid)
	}

	filter := bson.M{
		"accommodationId": bson.M{"$in": accommodationOids}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {

		return nil, err
	}

	var availabilities model.Availabilities
	err = cursor.All(ctx, &availabilities)
	if err != nil {

		return nil, err
	}

	hostIds := make([]string, 0)
	for _, val := range availabilities {
		exists := false
		for _, id := range hostIds {
			if id == val.HostId {
				exists = true
				break
			}
		}

		// If the hostId doesn't exist, append it to hostIds
		if !exists {
			hostIds = append(hostIds, val.HostId)
		}
	}

	return hostIds, nil
}

func prominentHostHttp(hostId string) (bool, error) {
	client := &http.Client{}

	apiHost := os.Getenv("API_GATEWAY_HOST")
	req, err := http.NewRequest("GET", "http://"+apiHost+":8000/api-2/accommodation/prominent-host/"+hostId, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) == "true" {
		return true, nil
	} else {
		return false, nil
	}
}
