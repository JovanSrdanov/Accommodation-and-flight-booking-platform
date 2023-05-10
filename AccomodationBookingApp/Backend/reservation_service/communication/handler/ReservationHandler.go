package handler

import (
	reservation "common/proto/reservation_service/generated"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reservation_service/domain/service"
)

type ReservationHandler struct {
	reservation.UnimplementedReservationServiceServer
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (handler ReservationHandler) CreateAvailability(ctx context.Context, in *reservation.CreateAvailabilityRequest) (*reservation.CreateAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}
func (handler ReservationHandler) GetAllMy(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllMyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMy not implemented")
}
func (handler ReservationHandler) UpdatePriceAndDate(ctx context.Context, in *reservation.UpdateRequest) (*reservation.UpdateRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePriceAndDate not implemented")
}
func (handler ReservationHandler) CreateReservation(ctx context.Context, in *reservation.CreateReservationRequest) (*reservation.CreateReservationRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReservation not implemented")
}
func (handler ReservationHandler) GetAllPendingReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllPendingReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPendingReservations not implemented")
}
func (handler ReservationHandler) GetAllRejectedReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllRejectedReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRejectedReservations not implemented")
}
func (handler ReservationHandler) RejectReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectReservation not implemented")
}
func (handler ReservationHandler) AcceptReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptReservation not implemented")
}
func (handler ReservationHandler) CancelReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelReservation not implemented")
}
