package handler

import (
	reservation "common/proto/reservation_service/generated"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	mapper := NewReservationMapper()
	id, err := handler.reservationService.CreateAvailability(mapper.mapFromCreateAvailability(in))

	if err != nil {
		return nil, err
	}
	return &reservation.CreateAvailabilityResponse{
		Id: id.String(),
	}, nil
}
func (handler ReservationHandler) GetAllMy(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllMyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMy not implemented")
}
func (handler ReservationHandler) UpdatePriceAndDate(ctx context.Context, in *reservation.UpdateRequest) (*reservation.UpdateRequest, error) {
	mapper := NewReservationMapper()
	_, err := handler.reservationService.UpdatePriceAndDate(mapper.mapFromUpdatePriceAndDate(in))

	if err != nil {
		return nil, err
	}
	return &reservation.UpdateRequest{}, nil
}
func (handler ReservationHandler) CreateReservation(ctx context.Context, in *reservation.CreateReservationRequest) (*reservation.CreateReservationRequest, error) {
	mapper := NewReservationMapper()
	_, err := handler.reservationService.CreateReservation(mapper.mapFromCreateReservation(in))

	if err != nil {
		return nil, err
	}
	return &reservation.CreateReservationRequest{}, nil
}
func (handler ReservationHandler) GetAllPendingReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllPendingReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPendingReservations not implemented")
}
func (handler ReservationHandler) GetAllRejectedReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllRejectedReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRejectedReservations not implemented")
}
func (handler ReservationHandler) RejectReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.RejectReservation(id)

	if err != nil {
		return nil, err
	}
	return &reservation.RejectReservationResponse{
		Id: id.String(),
	}, nil
}
func (handler ReservationHandler) AcceptReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.AcceptReservation(id)

	if err != nil {
		return nil, err
	}
	return &reservation.RejectReservationResponse{
		Id: id.String(),
	}, nil
}
func (handler ReservationHandler) CancelReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.CancelReservation(id)

	if err != nil {
		return nil, err
	}
	return &reservation.RejectReservationResponse{
		Id: id.String(),
	}, nil
}
