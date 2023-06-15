package handler

import (
	"accommodation_service/communication"
	"accommodation_service/domain/service"
	accommodation "common/proto/accommodation_service/generated"
	reservation "common/proto/reservation_service/generated"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"

	"github.com/google/uuid"
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
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}
	id, err := handler.accommodationService.Create(mapper.mapFromCreateRequest(loggedInId.String(), in))

	if err != nil {
		return nil, err
	}

	log.Println(id.Hex())

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
	// get account credentials id from logged-in user
	/*loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}

	// get account credentials from acc cred microservice
	accCredClient := client.NewAccountCredentialsClient("authorization_service:8000")
	accCred, err := accCredClient.GetById(ctx, &authorization.GetByIdRequest{Id: loggedInId.String()})
	if err != nil {
		return nil, err
	}

	// get user info
	userInfoId, err := uuid.Parse(accCred.GetAccountCredentials().GetUserProfileId())
	if err != nil {
		return nil, err
	}

	userProfileMapper := NewUserProfileMapper()
	userUpdatedInfo, err := handler.accommodationService.Update(userInfoId, userProfileMapper.mapUpdateRequestToUpdateDto(req))

	return userProfileMapper.mapUpdateDtoToUpdateRequest(userUpdatedInfo), nil*/
	return &accommodation.UpdateRequest{}, nil
}

func (handler AccommodationHandler) GetById(ctx context.Context, in *accommodation.GetByIdRequest) (*accommodation.GetByIdResponse, error) {
	/*id, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, err
	}
	userProfile, err := handler.accommodationService.GetById(id)
	if err != nil {
		return nil, err
	}

	mapper := NewUserProfileMapper()
	*/
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
	loggedInId := ctx.Value("id")
	err := handler.accommodationService.DeleteByHostId(loggedInId.(uuid.UUID).String())
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
	loggedInId := ctx.Value("id")
	accommodations, err := handler.accommodationService.GetAllMy(loggedInId.(uuid.UUID).String())
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
