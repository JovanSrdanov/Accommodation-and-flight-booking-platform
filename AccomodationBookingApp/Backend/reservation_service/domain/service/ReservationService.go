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

func (service ReservationService) GetAllMy() (model.Availabilities, error) {
	return service.reservationRepo.GetAllMy()
}

func (service ReservationService) UpdatePriceAndDate(priceWithDate *model.UpdatePriceAndDate) (*model.UpdatePriceAndDate, error) {
	return service.reservationRepo.UpdatePriceAndDate(priceWithDate)
}

func (service ReservationService) GetAllRejectedReservations() (*model.Reservation, error) {
	return service.reservationRepo.GetAllRejectedReservations()
}

func (service ReservationService) GetAllPendingReservations() (*model.Reservation, error) {
	return service.reservationRepo.GetAllPendingReservations()
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
