package controller

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Create godoc
// @Tags User
// @Param user body model.User true "User"
// @Consume application/json
// @Produce application/json
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.SimpleResponse
// @Router /user [post]
func (controller *UserController) Create(ctx *gin.Context) {
	var user model.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	id, err := controller.userService.Create(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewCreatedResponse(id))
}

// GetAll godoc
// @Tags User
// @Produce application/json
// @Success 200 {array} model.User
// @Failure 500 {object} dto.SimpleResponse
// @Router /user [get]
func (controller *UserController) GetAll(ctx *gin.Context) {
	users, err := controller.userService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewSimpleResponse("Error while reading from database"))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetById godoc
// @Tags User
// @Param id path string true "User ID"
// @Produce application/json
// @Success 200 {object} model.User
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /user/{id} [get]
func (controller *UserController) GetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	user, err := controller.userService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewSimpleResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete godoc
// @Tags User
// @Param id path string true "User ID"
// @Produce application/json
// @Success 200 {object} dto.SimpleResponse
// @Failure 400 {object} dto.SimpleResponse
// @Failure 404 {object} dto.SimpleResponse
// @Router /user/{id} [delete]
func (controller *UserController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewSimpleResponse(err.Error()))
		return
	}

	err = controller.userService.Delete(id)

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