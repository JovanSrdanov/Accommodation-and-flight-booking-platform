package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	authorization "FlightBookingApp/middleware"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	token "FlightBookingApp/token"
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

	response, err := controller.accountService.Register(registrationInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
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

	ctx.JSON(http.StatusOK, gin.H{"ACCESS TOKEN":accessTokenString, "REFRESH TOKEN":refreshTokenString})
}

func (controller *AccountController) RefreshAccessToken(ctx *gin.Context) {
	// TODO Stefan: endpoint should have :token at the end
	refreshToken := ctx.Param("token")
	
	// refresh token validation
	valid, claims := token.VerifyToken(refreshToken)
	if !valid || claims.TokenType != "refresh" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Invalid refresh token"})
		return
	}

	account, err := controller.accountService.GetById(claims.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"Error while generating the access token - user not found"})
		return
	}

	accessToken, err1 := token.GenerateAccessToken(account)
	if err1 != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"Error while generating the access token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"new access token":accessToken})
}

// GetAll godoc
// @Tags Account
// @Produce application/json
// @Success 200 {array} model.Account
// @Failure 500 {object} dto.SimpleResponse
// @Router /account [get]
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

	// getting the id from the token
	if len(ctx.Keys) == 0 {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"not authenticated"})
		}

	userID := ctx.Keys["ID"]
	userRole := ctx.Keys["Roles"]
	if userID == nil || userRole == nil{
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"can't get the user ID or roles"})
	}

	// a user can only see his information, unless he is an admin
	if id != userID && !authorization.RoleMatches(model.ADMIN, userRole.([]model.Role)) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"unauthorized access atempt"})
		return
	}

	account, err := controller.accountService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Delete godoc
// @Tags Account
// @Param id path string true "Account ID"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /account/{id} [delete]
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
