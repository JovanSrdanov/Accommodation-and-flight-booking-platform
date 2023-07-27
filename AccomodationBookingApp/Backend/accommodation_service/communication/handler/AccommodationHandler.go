package handler

import (
	"accommodation_service/communication"
	"accommodation_service/domain/service"
	"accommodation_service/utils"
	accommodation "common/proto/accommodation_service/generated"
	reservation "common/proto/reservation_service/generated"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationHandler struct {
	accommodation.UnimplementedAccommodationServiceServer
	accommodationService      service.AccommodationService
	reservationServiceAddress string
}

func NewAccommodationHandler(accommodationService service.AccommodationService, reservationServiceAddress string) *AccommodationHandler {
	return &AccommodationHandler{
		accommodationService:      accommodationService,
		reservationServiceAddress: reservationServiceAddress,
	}
}

func (handler AccommodationHandler) Create(ctx context.Context, in *accommodation.CreateRequest) (*accommodation.CreateResponse, error) {
	mapper := NewAccommodationMapper()
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	id, err := handler.accommodationService.Create(mapper.mapFromCreateRequest(loggedInId.String(), in))

	if err != nil {
		return nil, err
	}

	// Create availability base
	//Ovo samo ne radi jbg
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	_, err = reservationClient.CreateAvailabilityBase(ctx, &reservation.CreateAvailabilityBaseRequest{
		ReservationBase: &reservation.AvailabilityBase{
			AccommodationId:        id.Hex(),
			HostId:                 loggedInId.String(),
			IsAutomaticReservation: in.Accommodation.IsAutomaticReservation,
		},
	})
	if err != nil {
		return nil, err
	}

	return &accommodation.CreateResponse{
		Id: id.Hex(),
	}, nil
}

func (handler AccommodationHandler) Update(ctx context.Context, req *accommodation.UpdateRequest) (*accommodation.UpdateRequest, error) {
	return &accommodation.UpdateRequest{}, nil
}

func (handler AccommodationHandler) GetById(ctx context.Context, in *accommodation.GetByIdRequest) (*accommodation.GetByIdResponse, error) {
	mapper := NewAccommodationMapper()
	id, _ := primitive.ObjectIDFromHex(in.Id)

	res, err := handler.accommodationService.GetById(id)
	if err != nil {
		return nil, err
	}

	return mapper.mapToGetByIdResponse(res), nil
}

func (handler AccommodationHandler) Delete(ctx context.Context, in *accommodation.DeleteRequest) (*accommodation.DeleteResponse, error) {
	/*id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	err = handler.accommodationService.DeleteUserProfile(id)
	*/
	return &accommodation.DeleteResponse{}, nil
}

func (handler AccommodationHandler) DeleteByHostId(ctx context.Context, in *accommodation.EmptyRequest) (*accommodation.DeleteResponse, error) {
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	err = handler.accommodationService.DeleteByHostId(loggedInId.String())
	if err != nil {
		return &accommodation.DeleteResponse{}, err
	}
	return &accommodation.DeleteResponse{}, nil
}

func (handler AccommodationHandler) GetAll(ctx context.Context, in *accommodation.EmptyRequest) (*accommodation.GetAllResponse, error) {
	accommodations, err := handler.accommodationService.GetAll()
	if err != nil {
		return nil, err
	}

	mapper := NewAccommodationMapper()

	return mapper.mapToGetAllResponse(accommodations), nil
}

func (handler AccommodationHandler) GetAllMy(ctx context.Context, in *accommodation.GetMyRequest) (*accommodation.GetAllResponse, error) {
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	accommodations, err := handler.accommodationService.GetAllMy(loggedInId.String())
	if err != nil {
		return nil, err
	}

	mapper := NewAccommodationMapper()

	return mapper.mapToGetAllResponse(accommodations), nil
}

func (handler AccommodationHandler) GetAmenities(ctx context.Context, in *accommodation.EmptyRequest) (*accommodation.GetAmenitiesResponse, error) {
	amenities, err := handler.accommodationService.GetAmenities()
	if err != nil {
		return nil, err
	}

	return &accommodation.GetAmenitiesResponse{
		Amenities: amenities,
	}, nil
}

func (handler AccommodationHandler) SearchAccommodation(ctx context.Context, in *accommodation.SearchRequest) (*accommodation.GetAllResponse, error) {
	mapper := NewAccommodationMapper()

	accommodations, err := handler.accommodationService.SearchAccommodation(mapper.mapFromSearchRequest(in))
	if err != nil {
		return nil, err
	}

	return mapper.mapToGetAllResponse(accommodations), nil
}
