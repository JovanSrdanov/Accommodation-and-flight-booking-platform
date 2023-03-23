package controller

import (
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) AccountController {
	return AccountController{
		accountService: accountService,
	}
}

func (controller *AccountController) Register(ctx *gin.Context) {
	var registrationInfo dto.AccountRegistration

	if err := ctx.ShouldBindJSON(&registrationInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	hashedPassword, err := utils.HashPassword(registrationInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse(err.Error()))
		return
	}

	newAccount := model.Account {
		ID: uuid.New(),
		Username: registrationInfo.Username,
		Password: hashedPassword,
		Email: registrationInfo.Email,
		Role: model.REGULAR_USER,
		IsActivated: false,
	}

	account, err := controller.accountService.Register(newAccount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
	}

	//not sending hashed password in the reponse
	response := dto.CreateUserResponse {
		ID: account.ID,
		Username: account.Username,
		Email: account.Email,
		Role: account.Role,
		IsActivated: account.IsActivated,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (controller *AccountController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.accountService.GetAll())
}

func (controller *AccountController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	account, err := controller.accountService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
	}

	ctx.JSON(http.StatusOK, account)
}

func (controller *AccountController) Delete(ctx *gin.Context) {
	//TODO Stefan
	panic("not implemented")
}