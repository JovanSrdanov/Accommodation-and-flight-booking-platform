package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"FlightBookingApp/utils"
	"FlightBookingApp/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type FlightController struct {
	flightService service.FlightService
}

var validate *validator.Validate

func NewFlightController(flightService service.FlightService) *FlightController {
	validate = validator.New()
	validate.RegisterValidation("not-before-current-date", validators.NotBeforeCurrentDate)
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
		U tom slucaju morao bi se napraviti mini parser za ove generic poruke
	*/
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}
	err = validate.Struct(&flight)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = flight.SelfValidate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	id, err := controller.flightService.Create(&flight)
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

// Cancel godoc
// @Tags Flight
// @Param id path string true "Flight ID"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /flight/{id} [patch]
func (controller *FlightController) Cancel(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = controller.flightService.Cancel(id)

	if err != nil {
		switch err.(type) {
		case errors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
			return
		case errors.FlightPassedError:
			ctx.JSON(http.StatusConflict, dto.NewSimpleResponse(err.Error()))
			return
		default:
			ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, dto.NewSimpleResponse("Flight canceled"))
}

// Search godoc
// @Tags Flight
// @Produce application/json
// @Success 200 {array} utils.Page
// @Failure 400 {object} dto.SimpleResponse
// @Failure 500 {object} dto.SimpleResponse
// @Param departureDate query string true "Departure date, must be in this format YYYY-MM-DD" format(yyyy-mm-dd)
// @Param destinationCountry query string true "Destination country"
// @Param destinationCity query string true "Destination city"
// @Param startPointCountry query string true "Starting point country"
// @Param startPointCity query string true "Starting point  city"
// @Param desiredNumberOfSeats query int true "Desired number of seats"
// @Param pageNumber query int true "Page number"
// @Param resultsPerPage query int true "Results per page"
// @Param sortDirection query string true "Sort direction (asc, dsc, no_sort)"
// @Param sortType query string true "Sort type, it can be departureDateTime or price"
// @Router /flight/search [get]
func (controller *FlightController) Search(ctx *gin.Context) {
	flightSearchParameters, err := dto.NewFlightSearchParameters(
		ctx.Query("departureDate"),
		ctx.Query("destinationCountry"),
		ctx.Query("destinationCity"),
		ctx.Query("startPointCountry"),
		ctx.Query("startPointCity"),
		ctx.Query("desiredNumberOfSeats"),
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	pageInfo, err := utils.NewPageInfo(ctx.Query("pageNumber"), ctx.Query("resultsPerPage"), ctx.Query("sortDirection"),
		ctx.Query("sortType"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	flightsPage, err := controller.flightService.Search(flightSearchParameters, pageInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, flightsPage)
}
