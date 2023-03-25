package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type FlightController struct {
	flightService service.FlightService
}

func NewFlightController(flightService service.FlightService) *FlightController {
	return &FlightController{
		flightService: flightService,
	}
}

// Create godoc
// @Tags Flight
// @Param flight body model.Flight true "Flight"
// @Consume application/json
// @Produce application/json
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.SimpleResponse
// @Router /flight [post]
func (controller *FlightController) Create(ctx *gin.Context) {
	var flight model.Flight

	err := ctx.ShouldBindJSON(&flight)
	/*
		TODO Aleksandar: da li da pravimo custom message za neuspeli binding?
		U tom slucaju morao bi se napraviti mini parser za ove generic poruke
	*/
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	//Service call and return
	id, err := controller.flightService.Create(flight)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewCreatedResponse(id))
}

// GetAll godoc
// @Tags Flight
// @Produce application/json
// @Success 200 {array} model.Flight
// @Failure 500 {object} dto.SimpleResponse
// @Router /flight [get]
func (controller *FlightController) GetAll(ctx *gin.Context) {
	flights, err := controller.flightService.GetAll()

	if err != nil {
		//Couldn't connect to database
		//TODO Aleksandar: koji status code?
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, flights)
}

// GetById godoc
// @Tags Flight
// @Param id path string true "Flight ID"
// @Produce application/json
// @Success 200 {object} model.Flight
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /flight/{id} [get]
func (controller *FlightController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	flight, err := controller.flightService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

// Delete godoc
// @Tags Flight
// @Param id path string true "Flight ID"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /flight/{id} [delete]
func (controller *FlightController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = controller.flightService.Delete(id)

	if err != nil {
		switch err.(type) {
		case errors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
			return
		default:
			ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, dto.NewSimpleResponse("Entity deleted"))
}
