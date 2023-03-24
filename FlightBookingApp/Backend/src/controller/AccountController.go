package controller

import (
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO Stefan: podesi swagger

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{
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
		Username: registrationInfo.Username,
		Password: hashedPassword,
		Email: registrationInfo.Email,
		Role: model.REGULAR_USER,
		IsActivated: false,
	}

	id, err := controller.accountService.Register(newAccount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
	}

	newAccount.ID = id

	//not returning sensitive information
	response := dto.CreateUserResponse {
		ID: newAccount.ID,
		Role: newAccount.Role,
		IsActivated: newAccount.IsActivated,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (controller *AccountController) GetAll(ctx *gin.Context) {
	accounts, err := controller.accountService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (controller *AccountController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	account, err := controller.accountService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (controller *AccountController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = controller.accountService.Delete(id)

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