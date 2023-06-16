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

type TicketController struct {
	ticketService service.TicketService
	jwtService    service.JwtService
	apiKeyService *service.ApiKeyService
}

func NewTicketController(ticketService service.TicketService, jwtService service.JwtService, apiKeyService *service.ApiKeyService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
		jwtService:    jwtService,
		apiKeyService: apiKeyService,
	}
}

// Create godoc
// @Tags Ticket
// @Param ticket body model.Ticket true "Ticket"
// @Consume application/json
// @Produce application/json
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.SimpleResponse
// @Router /ticket [post]
func (controller *TicketController) Create(ctx *gin.Context) {
	var ticket model.Ticket

	err := ctx.ShouldBindJSON(&ticket)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	id, err := controller.ticketService.Create(ticket)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewCreatedResponse(id))
}

// GetAll godoc
// @Security bearerAuth
// @Tags Ticket
// @Produce application/json
// @Success 200 {array} model.Ticket
// @Failure 500 {object} dto.SimpleResponse
// @Router /ticket [get]
func (controller *TicketController) GetAll(ctx *gin.Context) {
	tickets, err := controller.ticketService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

// GetById godoc
// @Tags Ticket
// @Param id path string true "Ticket ID"
// @Produce application/json
// @Success 200 {object} model.Ticket
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /ticket/{id} [get]
func (controller *TicketController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ticket, err := controller.ticketService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}

// Delete godoc
// @Tags Ticket
// @Param id path string true "Ticket ID"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /ticket/{id} [delete]
func (controller *TicketController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = controller.ticketService.Delete(id)

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

// BuyTicket godoc
// @Security bearerAuth
// @Tags Ticket
// @Param ticket body dto.BuyTicketDto true "BuyTicketDto"
// @Consume application/json
// @Produce application/json
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.SimpleResponse
// @Router /ticket/buy [post]
func (controller *TicketController) BuyTicket(ctx *gin.Context) {
	//DTO i njegova validacija
	var buyTicketDto dto.BuyTicketDto

	err := ctx.ShouldBindJSON(&buyTicketDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	//Izvlacenje informacija o ulogovanom korisniku iz jwt-a
	accountID := ctx.Keys["ID"]
	user, err := controller.jwtService.GetUser(accountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse(err.Error()))
		return
	}

	//TODO Strahinja: buyer treba iz API kljuca
	ticket := model.NewTicket(user, user, buyTicketDto.FlightId)

	//Kreiranje tiketa u bazi
	id, err := controller.ticketService.BuyTicket(ticket, buyTicketDto.NumberOfTickets)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewCreatedResponse(id))
}

func (controller *TicketController) BuyTicketApiKey(ctx *gin.Context) {
	//DTO i njegova validacija
	var buyTicketApiKeyDto dto.BuyTicketApiKeyDto

	err := ctx.ShouldBindJSON(&buyTicketApiKeyDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	apiKey, err := controller.apiKeyService.GetByValue(buyTicketApiKeyDto.ApiKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	user, err := controller.jwtService.GetUser(apiKey.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse(err.Error()))
		return
	}

	ticket := model.NewTicket(user, user, buyTicketApiKeyDto.FlightId)

	//Kreiranje tiketa u bazi
	id, err := controller.ticketService.BuyTicket(ticket, buyTicketApiKeyDto.NumberOfTickets)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewCreatedResponse(id))
}

// GetAllForCustomer godoc
// @Security bearerAuth
// @Tags Ticket
// @Produce application/json
// @Success 200 {array} model.Ticket
// @Failure 500 {object} dto.SimpleResponse
// @Router /ticket/myTickets [get]
func (controller *TicketController) GetAllForCustomer(ctx *gin.Context) {
	accountID := ctx.Keys["ID"]
	user, err := controller.jwtService.GetUser(accountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse(err.Error()))
		return
	}

	tickets, err := controller.ticketService.GetAllForUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}
