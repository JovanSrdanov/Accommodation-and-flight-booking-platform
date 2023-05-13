package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reservation_service/domain/model"
	"reservation_service/domain/repository"
)

type ReservationService struct {
	reservationRepo repository.IReservationRepository
}

func NewReservationService(reservationRepo repository.IReservationRepository) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo}
}

func (service ReservationService) CreateAvailability(availability *model.AvailabilityRequest) (primitive.ObjectID, error) {
	return service.reservationRepo.CreateAvailability(availability)
}

func (service ReservationService) GetAllMy(hostId string) (model.Availabilities, error) {
	return service.reservationRepo.GetAllMy(hostId)
}

func (service ReservationService) UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error) {
	return service.reservationRepo.UpdatePriceAndDate(priceWithDate)
}

func (service ReservationService) GetAllAcceptedReservations(hostId string) (model.Reservations, error) {
	return service.reservationRepo.GetAllAcceptedReservations(hostId)
}

func (service ReservationService) GetAllPendingReservations(hostId string) (model.Reservations, error) {
	return service.reservationRepo.GetAllPendingReservations(hostId)
}

func (service ReservationService) GetAllReservationsForGuest(guestId string) (model.Reservations, error) {
	return service.reservationRepo.GetAllReservationsForGuest(guestId)
}

func (service ReservationService) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	return service.reservationRepo.CreateReservation(reservation)
}

func (service ReservationService) RejectReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return service.reservationRepo.RejectReservation(id)
}

func (service ReservationService) AcceptReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return service.reservationRepo.AcceptReservation(id)
}

func (service ReservationService) CancelReservation(id primitive.ObjectID) (primitive.ObjectID, error) {
	return service.reservationRepo.CancelReservation(id)
}

func (service ReservationService) CreateAvailabilityBase(base *model.Availability) (primitive.ObjectID, error) {
	return service.reservationRepo.CreateAvailabilityBase(base)
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
