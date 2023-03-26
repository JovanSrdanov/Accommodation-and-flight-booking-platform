package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"FlightBookingApp/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	newAccount := model.Account{
		Username:    registrationInfo.Username,
		Password:    hashedPassword,
		Email:       registrationInfo.Email,
		Role:        model.REGULAR_USER,
		IsActivated: false,
	}

	id, err := controller.accountService.Register(newAccount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	newAccount.ID = id

	response := dto.CreateUserResponse{
		ID:          newAccount.ID,
		Role:        newAccount.Role,
		IsActivated: newAccount.IsActivated,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (controller *AccountController) VerifyEmail(ctx *gin.Context) {
	var account model.Account
	account.Username = ctx.Param("username")
	linkVerPass := ctx.Param("verPass")

	account, err := controller.accountService.GetByUsername(account.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("error geting verification hash in db by username"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.EmailVerificationHash), []byte(linkVerPass))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.NewSimpleResponse("unauthorized access"))
	}

	if time.Now().Local().After(account.VerificationTimeout) {
		ctx.JSON(http.StatusUnauthorized, dto.NewSimpleResponse("email verification link has expired, please try registering again"))
		return
	}

	account.IsActivated = true
	controller.accountService.Save(account)
	ctx.JSON(http.StatusOK, dto.NewSimpleResponse("email successfuly verified"))
}

func (controller *AccountController) Login(ctx *gin.Context) {
	var loginData dto.LoginRequest

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse("invalid request"))
	}

	accessTokenString, refreshTokenString, err := controller.accountService.Login(loginData)

	//TODO Stefan: fix error handleing

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("-------ACCESS TOKEN-------: %v-------REFRESH TOKEN-------: %v", accessTokenString, refreshTokenString))
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
