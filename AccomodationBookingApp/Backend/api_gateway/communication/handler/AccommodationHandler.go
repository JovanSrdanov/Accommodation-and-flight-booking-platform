package handler

import (
	"api_gateway/communication"
	"api_gateway/dto"
	accommodation "common/proto/accommodation_service/generated"
	reservation "common/proto/reservation_service/generated"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type AccommodationHandler struct {
	accommodationServiceAddress string
	reservationServiceAddress   string
}

func NewAccommodationHandler(accommodationServiceAddress string, reservationServiceAddress string) *AccommodationHandler {
	return &AccommodationHandler{
		accommodationServiceAddress: accommodationServiceAddress,
		reservationServiceAddress:   reservationServiceAddress,
	}
}

func (handler AccommodationHandler) Init(router *gin.RouterGroup) {
	userGroup := router.Group("/accommodation")
	userGroup.POST("/search", handler.SearchAccommodation)
}

func (handler AccommodationHandler) SearchAccommodation(ctx *gin.Context) {
	var searchDto dto.SearchAccommodationDto

	err := ctx.ShouldBindJSON(&searchDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	var firstRoundDto dto.SearchResponseDto
	err = handler.FindAccommodations(searchDto, firstRoundDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	log.Println(firstRoundDto)

	finalDto, err := handler.FindReservations(searchDto, firstRoundDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, finalDto)

}

func (handler AccommodationHandler) FindAccommodations(searchDto dto.SearchAccommodationDto, responseDto dto.SearchResponseDto) error {
	client := communication.NewAccommodationClient(handler.accommodationServiceAddress)
	response, err := client.SearchAccommodation(context.TODO(), &accommodation.SearchRequest{Filter: &accommodation.Filter{
		Amenities: searchDto.Amenities,
		Location:  searchDto.Location,
		MinGuests: searchDto.MinGuests,
		HostId:    "xxx",
	}})

	if err != nil {
		return err
	}

	for _, value := range response.Accommodation {

		id, err := primitive.ObjectIDFromHex(value.Id)
		if err != nil {
			return err
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

	return nil
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
			if oldAcc.ID == id {
				oldAcc.Price = foundIds.Price
				finalDto = append(finalDto, oldAcc)
			}
		}
	}

	return finalDto, nil
}
