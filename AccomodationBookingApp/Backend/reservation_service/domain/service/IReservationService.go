package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reservation_service/domain/model"
)

type IAvailabilityService interface {
	CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error)
	GetAllMy(hostId string) (model.Availabilities, error)
	UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error)
	CreateReservation(reservation *model.Reservation) (*model.Reservation, error)
	GetAllPendingReservations(hostId string) (model.Reservations, error)
	GetAllAcceptedReservations(hostId string) (model.Reservations, error)
	RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CreateAvailabilityBase(base *model.Availability) (primitive.ObjectID, error)
	GetAllReservationsForGuest(guestId string) (model.Reservations, error)
	SearchAccommodation(accommodationIds []*primitive.ObjectID, dateRange model.DateRange, numberOfGuests int32) ([]*model.SearchResponseDto, error)
}
