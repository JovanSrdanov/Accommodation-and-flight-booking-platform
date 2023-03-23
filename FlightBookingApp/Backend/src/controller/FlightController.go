package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FlightController struct {
	flightService service.FlightService
}

func NewFlightController(flightService service.FlightService) FlightController {
	return FlightController{
		flightService: flightService,
	}
}

func (controller *FlightController) Create(ctx *gin.Context) {
	var flight model.Flight
	//Map and validate
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
	ctx.JSON(http.StatusCreated, controller.flightService.Create(flight))
}
func (controller *FlightController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.flightService.GetAll())
}

func (controller *FlightController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

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

func (controller *FlightController) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
