package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	authorization "FlightBookingApp/middleware"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	token "FlightBookingApp/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// Register godoc
// @Tags Account
// @Param registrationInfo body dto.AccountRegistration true "Registration info"
// @Consume application/json
// @Produce application/json
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.SimpleResponse
// @Router /account/register [post]
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

// VerifyEmail godoc
// @Tags Account
// @Param username path string true "Username"
// @Param verPass path string true "Email verification password"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 500 {object} dto.SimpleResponse
// @Failure 401 {object} dto.SimpleResponse
// @Router /account/emailver/{username}/{verPass} [get]
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
	//ctx.JSON(http.StatusOK, dto.NewSimpleResponse("email successfuly verified"))
	ctx.Redirect(302, "http://localhost:3000/")
}

// Login godoc
// @Tags Account
// @Param loginData body dto.LoginRequest true "Login data"
// @Consume application/json
// @Produce application/json
// @Success 200 {object} string
// @Failure 400 {object} dto.SimpleResponse
// @Router /account/login [post]
func (controller *AccountController) Login(ctx *gin.Context) {
	var loginData dto.LoginRequest

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse("invalid request"))
	}

	accessTokenString, refreshTokenString, err := controller.accountService.Login(loginData)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.NewSimpleResponse(err.Error()))
		return
	}

	response := dto.LoginResponse{AccessToken: accessTokenString}

	//access tokes is returned to the client, but the refresh token is saved in a http-only cookie
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",              // not sure what this is
		MaxAge:   60 * 60 * 24 * 7, // 1 week
	}

	// TODO Stefan: save the refresh token in the db

	http.SetCookie(ctx.Writer, cookie)

	ctx.JSON(http.StatusOK, response)
}

// RefreshAccessToken godoc
// @Security bearerAuth
// @Tags Account
// @Param token path string true "Refresh token"
// @Produce application/json
// @Success 200 {object} string
// @Failure 401 {object} string "Invalid refresh token"
// @Failure 500 {object} string "Error while generating the token"
// @Router /account/refresh-token/{token} [get]
func (controller *AccountController) RefreshAccessToken(ctx *gin.Context) {
	//refreshToken := ctx.Param("token")
	refreshCookie, err := ctx.Request.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No cookie provided"})
		return
	}

	refreshToken := refreshCookie.Value

	account, err := controller.accountService.GetByRefreshToken(refreshToken)

	fmt.Print(account)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden,
			gin.H{"error": err.Error()})
		return
	}

	// refresh token validation
	err, claims := token.VerifyToken(refreshToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if claims.TokenType != "refresh" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not a refresh token"})
		return
	}

	accessToken, err1 := token.GenerateAccessToken(account)
	if err1 != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while generating the access token"})
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetAll godoc
// @Security bearerAuth
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

// GetById godoc
// @Security bearerAuth
// @Tags Account
// @Param id path string true "Account ID"
// @Produce application/json
// @Success 200 {object} model.Account
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Failure 500 {object} string "can't get the account ID or roles"
// @Failure 401 {object} string "unauthorized access atempt"
// @Router /account/{id} [get]
func (controller *AccountController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	// getting the id from the token
	if len(ctx.Keys) == 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "not authenticated"})
		return
	}

	userID := ctx.Keys["ID"]
	userRole := ctx.Keys["Roles"]
	if userID == nil || userRole == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can't get the account ID or roles"})
		return
	}

	// a user can only see his information, unless he is an admin
	if id != userID && !authorization.RoleMatches(model.ADMIN, userRole.([]model.Role)) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access atempt"})
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
// @Security bearerAuth
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
