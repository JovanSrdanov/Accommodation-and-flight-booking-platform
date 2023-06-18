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
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	userGroup.GET("/rating/:accommodationId", handler.GetRatingDetailForAccommodation)
	userGroup.GET("/rating/host/:hostId", handler.GetRatingDetailForHost)
	userGroup.GET("/recommend",
		middleware.ValidateToken(handler.tokenMaker),
		middleware.Authorization([]model.Role{model.Guest}),
		handler.GetRecommendedAccommodations)
}

func (handler AccommodationHandler) SearchAccommodation(ctx *gin.Context) {
	//ctxGrpc := createGrpcContextFromGinContext(ctx)

	var searchDto dto.SearchAccommodationDto

	err := ctx.ShouldBindJSON(&searchDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		return
	}

	var firstRoundDto dto.SearchResponseDto
	firstRoundDto, err = handler.FindAccommodations(searchDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		return
	}

	log.Println(firstRoundDto)

	secondRound, err := handler.FindReservations(searchDto, firstRoundDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		return
	}

	finalDto, err := handler.FindRating(searchDto, secondRound, context.TODO())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		return
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
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

		isProminent := false

		if searchDto.ProminentHost {
			isProminent, err = prominentHostHttp(val)
			if err != nil {
				return nil, err
			}
		}

		if ratingForAccommodation.Rating.AvgRating >= searchDto.MinRating &&
			val.Price <= searchDto.MaxPrice &&
			val.Price >= searchDto.MinPrice &&
			containsAll(val.Amenities, searchDto.Amenities) &&
			(!searchDto.ProminentHost || isProminent) {
			val.Rating = ratingForAccommodation.Rating.AvgRating
			responseSlice = append(responseSlice, val)
		}
	}

	return responseSlice, nil
}

func prominentHostHttp(val *dto.Accommodation) (bool, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8000/api-2/accommodation/prominent-host/"+val.HostId, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

func containsAll(slice1, slice2 []string) bool {
	if len(slice2) == 0 {
		return true
	}

	for _, value := range slice2 {
		found := false
		for _, item := range slice1 {
			if item == value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (handler AccommodationHandler) GetRatableAccommodations(ctx *gin.Context) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	accommodationClient := communication.NewAccommodationClient(handler.accommodationServiceAddress)
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)

	loggedInAccCredIdFromCtx := ctx.Keys["id"].(uuid.UUID).String()
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	protoResponse, err := reservationClient.GetAllRatableAccommodationsForGuest(ctxGrpc, &reservation.GuestIdRequest{GuestId: loggedInAccCredIdFromCtx})
	if err != nil {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable accommodations?"})
		return
	}

	dtoSlice := make([]*dto.Accommodation, 0)
	for _, accId := range protoResponse.AccommodationIds {
		accommodationProto, err2 := accommodationClient.GetById(ctxGrpc, &accommodation.GetByIdRequest{Id: accId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable accommodations?"})
			return
		}

		protoAccommodationRating, err2 := ratingClient.GetRatingGuestGaveAccommodation(ctxGrpc, &rating.GetRatingGuestGaveAccommodationRequest{
			AccommodationId: accId,
			GuestId:         loggedInAccCredIdFromCtx,
		})

		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
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
			Rating:    protoAccommodationRating.Rating,
		})
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, dtoSlice)
}

func (handler AccommodationHandler) GetRatableHosts(ctx *gin.Context) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)

	loggedInAccCredIdFromCtx := ctx.Keys["id"].(uuid.UUID).String()
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	protoHostIds, err := reservationClient.GetAllRatableHostsForGuest(ctxGrpc, &reservation.GuestIdRequest{GuestId: loggedInAccCredIdFromCtx})
	if err != nil {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtoSlice := make([]*dto.BasicHostInfo, 0)
	for _, hostId := range protoHostIds.HostIds {
		protoAccInfo, err2 := authorizationClient.GetById(ctxGrpc, &authorization.GetByIdRequest{Id: hostId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable host?"})
			return
		}

		protoUserInfo, err2 := userProfileClient.GetById(ctxGrpc, &user_profile.GetByIdRequest{Id: protoAccInfo.AccountCredentials.UserProfileId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		protoHostRating, err2 := ratingClient.GetRatingGuestGaveHost(ctxGrpc, &rating.GetRatingGuestGaveHostRequest{
			HostId:  hostId,
			GuestId: loggedInAccCredIdFromCtx,
		})

		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		dtoSlice = append(dtoSlice, &dto.BasicHostInfo{
			Username: protoAccInfo.AccountCredentials.Username,
			Name:     protoUserInfo.UserProfile.Name,
			Surname:  protoUserInfo.UserProfile.Surname,
			Email:    protoUserInfo.UserProfile.Email,
			HostId:   hostId,
			Rating:   protoHostRating.Rating,
		})
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, dtoSlice)
}

func (handler AccommodationHandler) IsHostProminent(ctx *gin.Context) {
	hostId := ctx.Param("hostId")

	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)

	ratingProto, err := ratingClient.CalculateRatingForHost(context.TODO(), &rating.RatingForHostRequest{HostId: hostId})
	if err != nil {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(ratingProto.Rating)

	if ratingProto.Rating.AvgRating <= 4.7 {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountSucc.Inc()
		log.Println("Counter incremented")
		ctx.JSON(http.StatusOK, false)
		return
	}

	protoAllReservations, err := reservationClient.GetAllReservationsForHost(ctx, &reservation.HostIdRequest{HostId: hostId})
	if err != nil {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountFail.Inc()
		log.Println("Counter incremented")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	numberOfAccepted := 0
	numberOfCanceled := 0
	durationOfReservedDays := 0

	for _, protoReservation := range protoAllReservations.Reservation {
		if protoReservation.Status == "accepted" {
			numberOfAccepted++

			startDate := time.Unix(protoReservation.DateRange.From, 0)
			endDate := time.Unix(protoReservation.DateRange.To, 0)

			duration := endDate.Sub(startDate)
			days := int(duration.Hours()/24) + 1
			durationOfReservedDays += days
		} else if protoReservation.Status == "canceled" {
			numberOfCanceled++
		}
	}
	cancelRate := 0.0
	if numberOfCanceled != 0 {
		cancelRate = float64(numberOfAccepted+numberOfCanceled) / float64(numberOfCanceled) * 100.0
	}

	if cancelRate > 5.0 || durationOfReservedDays < 50 {
		middleware.HttpReqCountTotal.Inc()
		middleware.HttpReqCountSucc.Inc()
		log.Println("Counter incremented")
		ctx.JSON(http.StatusOK, false)
		return
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, true)
}

func (handler AccommodationHandler) GetRatingDetailForAccommodation(ctx *gin.Context) {
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)

	accommodationId := ctx.Param("accommodationId")

	protoRatingDetails, err := ratingClient.GetRatingForAccommodation(ctx, &rating.RatingForAccommodationRequest{AccommodationId: accommodationId})
	if err != nil {
		return
	}

	guestsInfo := make([]*dto.AccommodationRatingRating, 0)

	for _, protoGuestInfo := range protoRatingDetails.Rating.Ratings {
		protoAccInfo, err2 := authorizationClient.GetById(ctx, &authorization.GetByIdRequest{Id: protoGuestInfo.GuestId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable host?"})
			return
		}

		protoUserInfo, err2 := userProfileClient.GetById(ctx, &user_profile.GetByIdRequest{Id: protoAccInfo.AccountCredentials.UserProfileId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		guestsInfo = append(guestsInfo, &dto.AccommodationRatingRating{
			GuestID: protoGuestInfo.GuestId,
			Date:    protoGuestInfo.Date,
			Rating:  protoGuestInfo.Rating,
			Name:    protoUserInfo.UserProfile.Name,
			Surname: protoUserInfo.UserProfile.Surname,
		})
	}

	response := dto.AccommodationRating{
		AvgRating:       protoRatingDetails.Rating.AvgRating,
		AccommodationID: protoRatingDetails.Rating.AccommodationId,
		Ratings:         guestsInfo,
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, response)
}

func (handler AccommodationHandler) GetRatingDetailForHost(ctx *gin.Context) {
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)

	hostId := ctx.Param("hostId")

	protoRatingDetails, err := ratingClient.GetRatingForHost(ctx, &rating.RatingForHostRequest{HostId: hostId})
	if err != nil {
		return
	}

	guestsInfo := make([]*dto.HostRatingRating, 0)

	for _, protoGuestInfo := range protoRatingDetails.Rating.Ratings {
		protoAccInfo, err2 := authorizationClient.GetById(ctx, &authorization.GetByIdRequest{Id: protoGuestInfo.GuestId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Big puc kod get ratable host?"})
			return
		}

		protoUserInfo, err2 := userProfileClient.GetById(ctx, &user_profile.GetByIdRequest{Id: protoAccInfo.AccountCredentials.UserProfileId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		guestsInfo = append(guestsInfo, &dto.HostRatingRating{
			GuestID: protoGuestInfo.GuestId,
			Date:    protoGuestInfo.Date,
			Rating:  protoGuestInfo.Rating,
			Name:    protoUserInfo.UserProfile.Name,
			Surname: protoUserInfo.UserProfile.Surname,
		})
	}

	response := dto.HostRating{
		AvgRating: protoRatingDetails.Rating.AvgRating,
		HostID:    protoRatingDetails.Rating.HostId,
		Ratings:   guestsInfo,
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, response)
}

func (handler AccommodationHandler) GetRecommendedAccommodations(ctx *gin.Context) {
	accommodationClient := communication.NewAccommodationClient(handler.accommodationServiceAddress)
	ratingClient := communication.NewRatingClient(handler.ratingServiceAddress)

	loggedInAccCredIdFromCtx := ctx.Keys["id"].(uuid.UUID).String()
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	protoRecommendations, err := ratingClient.GetRecommendedAccommodations(ctxGrpc, &rating.RecommendedAccommodationsRequest{GuestId: loggedInAccCredIdFromCtx})
	if err != nil {
		return
	}

	recommendations := make([]*dto.Recommendation, 0)

	for _, protoRecommendation := range protoRecommendations.Recommendation {
		accommodationInfoProto, err2 := accommodationClient.GetById(ctxGrpc, &accommodation.GetByIdRequest{Id: protoRecommendation.AccommodationId})
		if err2 != nil {
			middleware.HttpReqCountTotal.Inc()
			middleware.HttpReqCountFail.Inc()
			log.Println("Counter incremented")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}

		recommendations = append(recommendations, &dto.Recommendation{
			AccommodationID: protoRecommendation.AccommodationId,
			Rating:          protoRecommendation.Rating,
			Name:            accommodationInfoProto.Accommodation.Name,
			Address: dto.Address{
				Country:      accommodationInfoProto.Accommodation.Address.Country,
				City:         accommodationInfoProto.Accommodation.Address.City,
				Street:       accommodationInfoProto.Accommodation.Address.Street,
				StreetNumber: accommodationInfoProto.Accommodation.Address.StreetNumber,
			},
			MinGuests: accommodationInfoProto.Accommodation.MinGuests,
			MaxGuests: accommodationInfoProto.Accommodation.MaxGuests,
			Amenities: accommodationInfoProto.Accommodation.Amenities,
			Images:    accommodationInfoProto.Accommodation.Images,
			HostId:    accommodationInfoProto.Accommodation.HostId,
		})
	}

	middleware.HttpReqCountTotal.Inc()
	middleware.HttpReqCountSucc.Inc()
	log.Println("Counter incremented")
	ctx.JSON(http.StatusOK, recommendations)
}
