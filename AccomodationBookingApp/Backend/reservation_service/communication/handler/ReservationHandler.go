package handler

import (
	reservation "common/proto/reservation_service/generated"
	"context"
	"reservation_service/domain/service"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationHandler struct {
	reservation.UnimplementedReservationServiceServer
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

func (handler ReservationHandler) CreateAvailability(ctx context.Context, in *reservation.CreateAvailabilityRequest) (*reservation.CreateAvailabilityResponse, error) {
	mapper := NewReservationMapper()
	id, err := handler.reservationService.CreateAvailability(mapper.mapFromCreateAvailability(in))

	if err != nil {
		return nil, err
	}
	return &reservation.CreateAvailabilityResponse{
		Id: id.Hex(),
	}, nil
}
func (handler ReservationHandler) GetAllMy(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllMyResponse, error) {
	mapper := NewReservationMapper()

	loggedInId := ctx.Value("id")
	availabilities, err := handler.reservationService.GetAllMy(loggedInId.(uuid.UUID).String())
	if err != nil {
		return &reservation.GetAllMyResponse{}, err
	}

	return mapper.mapToGetAllMyResponse(availabilities), nil
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
	mappedReservation := mapper.mapFromCreateReservation(in)
	loggedInId := ctx.Value("id")
	mappedReservation.GuestId = loggedInId.(uuid.UUID).String()
	_, err := handler.reservationService.CreateReservation(mappedReservation)

	if err != nil {
		return nil, err
	}

	return &reservation.CreateReservationRequest{}, nil
}
func (handler ReservationHandler) GetAllPendingReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllPendingReservationsResponse, error) {
	mapper := NewReservationMapper()

	loggedInId := ctx.Value("id")
	reservations, numberOfCancellations, err := handler.reservationService.GetAllPendingReservations(loggedInId.(uuid.UUID).String())
	if err != nil {
		return &reservation.GetAllPendingReservationsResponse{}, err
	}

	return &reservation.GetAllPendingReservationsResponse{
		Reservation: mapper.mapToReservationsProtoDto(reservations, numberOfCancellations),
	}, nil
}
func (handler ReservationHandler) GetAllAcceptedReservations(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllAcceptedReservationsResponse, error) {
	mapper := NewReservationMapper()

	loggedInId := ctx.Value("id")
	reservations, err := handler.reservationService.GetAllAcceptedReservations(loggedInId.(uuid.UUID).String())
	if err != nil {
		return &reservation.GetAllAcceptedReservationsResponse{}, err
	}

	return &reservation.GetAllAcceptedReservationsResponse{
		Reservation: mapper.mapToReservationsProto(reservations),
	}, nil
}
func (handler ReservationHandler) GetAllReservationsForGuest(ctx context.Context, in *reservation.EmptyRequest) (*reservation.GetAllReservationsForGuestResponse, error) {
	mapper := NewReservationMapper()

	loggedInId := ctx.Value("id")
	reservations, err := handler.reservationService.GetAllReservationsForGuest(loggedInId.(uuid.UUID).String())
	if err != nil {
		return &reservation.GetAllReservationsForGuestResponse{}, err
	}

	return &reservation.GetAllReservationsForGuestResponse{
		Reservations: mapper.mapToReservationsProto(reservations),
	}, nil
}
func (handler ReservationHandler) RejectReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.RejectReservation(id)

	if err != nil {
		return nil, err
	}

	return &reservation.RejectReservationResponse{
		Id: id.Hex(),
	}, nil
}
func (handler ReservationHandler) AcceptReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.AcceptReservation(id)

	if err != nil {
		return nil, err
	}
	return &reservation.RejectReservationResponse{
		Id: id.Hex(),
	}, nil
}
func (handler ReservationHandler) CancelReservation(ctx context.Context, in *reservation.ChangeStatusRequest) (*reservation.RejectReservationResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err := handler.reservationService.CancelReservation(id)

	if err != nil {
		return nil, err
	}
	return &reservation.RejectReservationResponse{
		Id: id.Hex(),
	}, nil
}
func (handler ReservationHandler) CreateAvailabilityBase(ctx context.Context, in *reservation.CreateAvailabilityBaseRequest) (*reservation.EmptyRequest, error) {
	mapper := NewReservationMapper()
	_, err := handler.reservationService.CreateAvailabilityBase(mapper.mapFromCreateAvailabilityBase(in.ReservationBase.HostId, in))

	if err != nil {
		return nil, err
	}
	return &reservation.EmptyRequest{}, nil
}
func (handler ReservationHandler) GuestHasActiveReservations(ctx context.Context, in *reservation.GuestHasActiveReservationsRequest) (*reservation.GuestHasActiveReservationsResponse, error) {
	id, err := uuid.Parse(in.GuestId)
	if err != nil {
		return nil, err
	}
	hasActiveReservations, err := handler.reservationService.GuestHasActiveReservations(id)
	if err != nil {
		return nil, err
	}

	return &reservation.GuestHasActiveReservationsResponse{HasActiveReservations: hasActiveReservations}, nil
}
func (handler ReservationHandler) SearchAccommodation(ctx context.Context, in *reservation.SearchRequest) (*reservation.SearchResponse, error) {
	mapper := NewReservationMapper()

	searchResponse, err := handler.reservationService.SearchAccommodation(mapper.mapFromSearchRequest(in))
	if err != nil {
		return nil, err
	}

	return mapper.mapToSearchResponse(searchResponse), nil
}
func (handler ReservationHandler) HostHasActiveReservations(ctx context.Context, in *reservation.HostHasActiveReservationsRequest) (*reservation.HostHasActiveReservationsResponse, error) {

	id, err := uuid.Parse(in.HostId)
	if err != nil {
		return nil, err
	}

	activeReservations, err := handler.reservationService.GetAllAcceptedReservations(id.String())
	if err != nil {
		return nil, err
	}

	if len(activeReservations) < 1 {
		return &reservation.HostHasActiveReservationsResponse{HasActiveReservations: false}, nil
	}
	return &reservation.HostHasActiveReservationsResponse{HasActiveReservations: true}, nil
}
func (handler ReservationHandler) GetAllRatableAccommodationsForGuest(ctx context.Context, in *reservation.GuestIdRequest) (*reservation.AccommodationsIdsResponse, error) {
	res, err := handler.reservationService.GetAllRatableAccommodationsForGuest(in.GuestId)
	if err != nil {
		return nil, err
	}
	return &reservation.AccommodationsIdsResponse{AccommodationIds: res}, nil
}
func (handler ReservationHandler) GetAllRatableHostsForGuest(ctx context.Context, in *reservation.GuestIdRequest) (*reservation.HostIdsResponse, error) {
	res, err := handler.reservationService.GetAllRatableHostsForGuest(in.GuestId)
	if err != nil {
		return nil, err
	}
	return &reservation.HostIdsResponse{HostIds: res}, nil
}
func (handler ReservationHandler) GetAllReservationsForHost(ctx context.Context, in *reservation.HostIdRequest) (*reservation.GetAllReservationsResponse, error) {
	mapper := NewReservationMapper()
	protoReservations, err := handler.reservationService.GetAllReservationsForHost(in.HostId)
	if err != nil {
		return nil, err
	}

	return &reservation.GetAllReservationsResponse{
		Reservation: mapper.mapToReservationsProto(protoReservations),
	}, nil
}
