package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reservation_service/domain/model"
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
	GetAllRejectedReservations() (*model.Reservation, error)
	RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error)
*/

func (repo ReservationRepositoryMongo) CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error) {
	return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) GetAllRejectedReservations() (*model.Reservation, error) {
	return &model.Reservation{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) GetAllPendingReservations() (*model.Reservation, error) {
	return &model.Reservation{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	return &model.Reservation{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error) {
	return &model.UpdatePriceAndDate{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) GetAllMy() (model.Availabilities, error) {
	return model.Availabilities{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}

func (repo ReservationRepositoryMongo) AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return primitive.ObjectID{}, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}
