package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type AirportController struct {
	airportService service.AirportService
}

func NewAirportController(airportService service.AirportService) *AirportController {
	return &AirportController{
		airportService: airportService,
	}
}

// GetAll godoc
// @Tags Airport
// @Produce application/json
// @Success 200 {array} model.Airport
// @Failure 500 {object} dto.SimpleResponse
// @Router /airport [get]
func (controller *AirportController) GetAll(ctx *gin.Context) {
	airports, err := controller.airportService.GetAll()

	if err != nil {
		//Couldn't connect to database
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, airports)
}

// GetById godoc
// @Tags Airport
// @Param id path string true "Airport ID"
// @Produce application/json
// @Success 200 {object} model.Airport
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /airport/{id} [get]
func (controller *AirportController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	airport, err := controller.airportService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, airport)
}
