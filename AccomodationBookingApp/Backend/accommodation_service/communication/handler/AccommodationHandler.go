package handler

import (
	"accommodation_service/domain/service"
	accommodation "common/proto/accommodation_service/generated"
	"context"
)

type AccommodationHandler struct {
	accommodation.UnimplementedAccommodationServiceServer
	accommodationService service.AccommodationService
}

func NewAccommodationHandler(accommodationService service.AccommodationService) *AccommodationHandler {
	return &AccommodationHandler{accommodationService: accommodationService}
}

func (handler AccommodationHandler) Create(ctx context.Context, in *accommodation.CreateRequest) (*accommodation.CreateResponse, error) {
	mapper := NewAccommodationMapper()
	id, err := handler.accommodationService.Create(mapper.mapFromCreateRequest(in))

	if err != nil {
		return nil, err
	}
	return &accommodation.CreateResponse{
		Id: id.String(),
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

	return &accommodation.GetByIdResponse{}, nil
}

func (handler AccommodationHandler) Delete(ctx context.Context, in *accommodation.DeleteRequest) (*accommodation.DeleteResponse, error) {
	/*id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	err = handler.accommodationService.Delete(id)
	*/
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
