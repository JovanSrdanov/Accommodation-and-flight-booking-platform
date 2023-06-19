package repository

import (
	"reservation_service/domain/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IReservationRepository interface {
	CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error)
	GetAllMy(hostId string) (model.Availabilities, error)
	UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error)
	CreateReservation(reservation *model.Reservation) (*model.Reservation, error)
	GetAllPendingReservations(hostId string) (model.Reservations, []int32, error)
	GetAllAcceptedReservations(hostId string) (model.Reservations, error)
	RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error)
	CreateAvailabilityBase(base *model.Availability) (primitive.ObjectID, error)
	GetAllReservationsForGuest(guestId string) (model.Reservations, error)
	GuestHasActiveReservations(guestID uuid.UUID) (bool, error)
	SearchAccommodation(accommodationIds []*primitive.ObjectID, dateRange model.DateRange, numberOfGuests int32) ([]*model.SearchResponseDto, error)
	DeleteAvailabilitiesAndReservationsByAccommodationId(accommodationId primitive.ObjectID) error

	GetAllReservationsForHost(hostId string) (model.Reservations, error)
	GetAllRatableAccommodationsForGuest(guestId string) ([]string, error)
	GetAllRatableHostsForGuest(guestId string) ([]string, error)
}
