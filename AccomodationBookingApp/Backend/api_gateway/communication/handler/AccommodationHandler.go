package handler

import (
	"api_gateway/communication"
	"api_gateway/communication/middleware"
	"api_gateway/dto"
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	accommodation "common/proto/accommodation_service/generated"
	authorization "common/proto/authorization_service/generated"
	rating "common/proto/rating_service/generated"
	reservation "common/proto/reservation_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type AccommodationHandler struct {
	accommodationServiceAddress string
	reservationServiceAddress   string
	authorizationServiceAddress string
	userProfileServiceAddress   string
	ratingServiceAddress        string
	tokenMaker                  token.Maker
}

func NewAccommodationHandler(accommodationServiceAddress string, reservationServiceAddress string,
	authorizationServiceAddress string, userProfileServiceAddress string, ratingServiceAddress string, tokenMaker token.Maker) *AccommodationHandler {
	return &AccommodationHandler{
		accommodationServiceAddress: accommodationServiceAddress,
		reservationServiceAddress:   reservationServiceAddress,
		authorizationServiceAddress: authorizationServiceAddress,
		userProfileServiceAddress:   userProfileServiceAddress,
		ratingServiceAddress:        ratingServiceAddress,
		tokenMaker:                  tokenMaker,
	}
}

func (handler AccommodationHandler) Init(router *gin.RouterGroup) {
	userGroup := router.Group("/accommodation")
	userGroup.POST("/search", handler.SearchAccommodation)
	userGroup.GET("/ratable/accommodations",
		middleware.ValidateToken(handler.tokenMaker),
		middleware.Authorization([]model.Role{model.Guest}),
		handler.GetRatableAccommodations)
	userGroup.GET("/ratable/hosts",
		middleware.ValidateToken(handler.tokenMaker),
		middleware.Authorization([]model.Role{model.Guest}),
		handler.GetRatableHosts)
	userGroup.GET("/prominent-host/:hostId", handler.IsHostProminent)
}

func (handler AccommodationHandler) SearchAccommodation(ctx *gin.Context) {
	//ctxGrpc := createGrpcContextFromGinContext(ctx)

	var searchDto dto.SearchAccommodationDto

	err := ctx.ShouldBindJSON(&searchDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	var firstRoundDto dto.SearchResponseDto
	firstRoundDto, err = handler.FindAccommodations(searchDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	log.Println(firstRoundDto)

	secondRound, err := handler.FindReservations(searchDto, firstRoundDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	finalDto, err := handler.FindRating(searchDto, secondRound, context.TODO())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, finalDto)

}

func (handler AccommodationHandler) FindAccommodations(searchDto dto.SearchAccommodationDto) (dto.SearchResponseDto, error) {
	var responseDto dto.SearchResponseDto

	client := communication.NewAccommodationClient(handler.accommodationServiceAddress)
	response, err := client.SearchAccommodation(context.TODO(), &accommodation.SearchRequest{Filter: &accommodation.Filter{
		Amenities: searchDto.Amenities,
		Location:  searchDto.Location,
		MinGuests: searchDto.MinGuests,
		HostId:    "xxx",
	}})

	if err != nil {
		return dto.SearchResponseDto{}, err
	}

	for _, value := range response.Accommodation {

		id, err := primitive.ObjectIDFromHex(value.Id)
		if err != nil {
			return dto.SearchResponseDto{}, err
		}

		responseDto = append(responseDto, &dto.Accommodation{
			ID:   id,
			Name: value.Name,
			Address: dto.Address{
				Country:      value.Address.Country,
				City:         value.Address.City,
				Street:       value.Address.Street,
				StreetNumber: value.Address.StreetNumber,
			},
			MinGuests: value.MinGuests,
			MaxGuests: value.MaxGuests,
			Amenities: value.Amenities,
			Images:    value.Images,
			HostId:    value.HostId,
			Price:     0,
		})
	}

	return responseDto, nil
}

func (handler AccommodationHandler) FindReservations(searchDto dto.SearchAccommodationDto, firstRoundDto dto.SearchResponseDto) (dto.SearchResponseDto, error) {
	client := communication.NewReservationClient(handler.reservationServiceAddress)

	accIds := make([]string, 0)

	for _, id := range firstRoundDto {
		accIds = append(accIds, id.ID.Hex())
	}

	res, err := client.SearchAccommodation(context.TODO(), &reservation.SearchRequest{Filter: &reservation.Filter{
		AccommodationIds: accIds,
		DateRange: &reservation.DateRange{
			From: searchDto.StartDate,
			To:   searchDto.EndDate,
		},
		NumberOfGuests: searchDto.MinGuests,
	}})

	if err != nil {
		return nil, err
	}

	var finalDto dto.SearchResponseDto

	for _, oldAcc := range firstRoundDto {
		for _, foundIds := range res.SearchResponse {
			id, err2 := primitive.ObjectIDFromHex(foundIds.AccommodationId)
			if err2 != nil {
				return nil, err
			}
			//TODO: Strahinja dodati ovde posle za max price proveru (&&)
			if oldAcc.ID == id {
				oldAcc.Price = foundIds.Price
				finalDto = append(finalDto, oldAcc)
			}
		}
	}

	return finalDto, nil
}

func (handler AccommodationHandler) FindRating(searchDto dto.SearchAccommodationDto, secondRoundDto dto.SearchResponseDto, ctx context.Context) (dto.SearchResponseDto, error) {
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)

	responseSlice := make([]*dto.Accommodation, 0)
	for _, val := range secondRoundDto {
		ratingForAccommodation, err := ratingClient.CalculateRatingForAccommodation(ctx, &rating.RatingForAccommodationRequest{AccommodationId: val.ID.Hex()})
		if err != nil {
			return nil, err
		}

		if ratingForAccommodation.Rating.AvgRating >= searchDto.MinRating {
			val.Rating = ratingForAccommodation.Rating.AvgRating
			responseSlice = append(responseSlice, val)
		}
	}

	return responseSlice, nil
}

func (handler AccommodationHandler) GetRatableAccommodations(ctx *gin.Context) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	accommodationClient := communication.NewAccommodationClient(handler.accommodationServiceAddress)

	loggedInAccCredIdFromCtx := ctx.Keys["id"].(uuid.UUID).String()
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	protoResponse, err := reservationClient.GetAllRatableAccommodationsForGuest(ctxGrpc, &reservation.GuestIdRequest{GuestId: loggedInAccCredIdFromCtx})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable accommodations?"})
		return
	}

	dtoSlice := make([]*dto.Accommodation, 0)
	for _, accId := range protoResponse.AccommodationIds {
		accommodationProto, err2 := accommodationClient.GetById(ctxGrpc, &accommodation.GetByIdRequest{Id: accId})
		if err2 != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable accommodations?"})
			return
		}
		id, _ := primitive.ObjectIDFromHex(accId)
		dtoSlice = append(dtoSlice, &dto.Accommodation{
			ID:   id,
			Name: accommodationProto.Accommodation.Name,
			Address: dto.Address{
				Country:      accommodationProto.Accommodation.Address.Country,
				City:         accommodationProto.Accommodation.Address.City,
				Street:       accommodationProto.Accommodation.Address.Street,
				StreetNumber: accommodationProto.Accommodation.Address.StreetNumber,
			},
			MinGuests: accommodationProto.Accommodation.MinGuests,
			MaxGuests: accommodationProto.Accommodation.MaxGuests,
			Amenities: accommodationProto.Accommodation.Amenities,
			Images:    accommodationProto.Accommodation.Images,
			HostId:    accommodationProto.Accommodation.HostId,
			Price:     -1,
		})
	}

	ctx.JSON(http.StatusOK, dtoSlice)
}

func (handler AccommodationHandler) GetRatableHosts(ctx *gin.Context) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)

	loggedInAccCredIdFromCtx := ctx.Keys["id"].(uuid.UUID).String()
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	protoHostIds, err := reservationClient.GetAllRatableHostsForGuest(ctxGrpc, &reservation.GuestIdRequest{GuestId: loggedInAccCredIdFromCtx})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtoSlice := make([]*dto.BasicUserInfo, 0)
	for _, hostId := range protoHostIds.HostIds {
		protoAccInfo, err2 := authorizationClient.GetById(ctxGrpc, &authorization.GetByIdRequest{Id: hostId})
		if err2 != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable host?"})
			return
		}

		protoUserInfo, err2 := userProfileClient.GetById(ctxGrpc, &user_profile.GetByIdRequest{Id: protoAccInfo.AccountCredentials.UserProfileId})
		if err2 != nil {
			return
		}

		dtoSlice = append(dtoSlice, &dto.BasicUserInfo{
			Username: protoAccInfo.AccountCredentials.Username,
			Name:     protoUserInfo.UserProfile.Name,
			Surname:  protoUserInfo.UserProfile.Surname,
			Email:    protoUserInfo.UserProfile.Email,
			HostId:   hostId,
		})
	}

	ctx.JSON(http.StatusOK, dtoSlice)
}

func (handler AccommodationHandler) IsHostProminent(ctx *gin.Context) {
	hostId := ctx.Param("hostId")
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)

	ratingProto, err := ratingClient.CalculateRatingForHost(ctxGrpc, &rating.RatingForHostRequest{HostId: hostId})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(ratingProto.Rating)

	ctx.JSON(http.StatusOK, ratingProto.Rating.AvgRating)
}
