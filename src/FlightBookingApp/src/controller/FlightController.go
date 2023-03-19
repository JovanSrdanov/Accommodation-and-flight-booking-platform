package controller

import (
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type flightController struct {
	flightService service.FlightService
}

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
	//Map
	err := ctx.ShouldBindJSON(&flight)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	//Validate
	//TODO treba u modelu da se podesi bind properties
	//Service call and return
	ctx.JSON(http.StatusCreated, controller.flightService.Create(flight))
}
func (controller *flightController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.flightService.GetAll())
}

func (controller *flightController) GetById(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (controller *flightController) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
