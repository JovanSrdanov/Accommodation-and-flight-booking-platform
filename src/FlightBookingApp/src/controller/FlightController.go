package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type flightController struct {
	flightService service.FlightService
}

// TODO da li ovo izmestati u zasebnu klasu ili ide gas?
type FlightController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewFlightController(flightService service.FlightService) FlightController {
	return &flightController{
		flightService: flightService,
	}
}
func (controller *flightController) Create(ctx *gin.Context) {
	var flight model.Flight
	//Map and validate
	err := ctx.ShouldBindJSON(&flight)
	/*
		TODO da li da pravimo custom message za neuspeli binding?
		U tom slucaju morao bi se napraviti mini parser za ove generic poruke
	*/
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	//Service call and return
	ctx.JSON(http.StatusCreated, controller.flightService.Create(flight))
}
func (controller *flightController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.flightService.GetAll())
}

func (controller *flightController) GetById(ctx *gin.Context) {
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

func (controller *flightController) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
